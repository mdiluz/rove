package internal

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
)

const (
	// PersistentData will allow the server to load and save it's state
	PersistentData = iota

	// EphemeralData will let the server neither load or save out any of it's data
	EphemeralData
)

// Server contains the relevant data to run a game server
type Server struct {

	// Internal state
	world *rove.World

	// Accountant
	accountant Accountant

	// gRPC server
	netListener net.Listener
	grpcServ    *grpc.Server

	// Config settings
	address        string
	persistence    int
	minutesPerTick int

	// sync point for sub-threads
	sync sync.WaitGroup

	// cron schedule for world ticks
	schedule *cron.Cron
}

// ServerOption defines a server creation option
type ServerOption func(s *Server)

// OptionAddress sets the server address for hosting
func OptionAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

// OptionPersistentData sets the server data to be persistent
func OptionPersistentData() ServerOption {
	return func(s *Server) {
		s.persistence = PersistentData
	}
}

// OptionTick defines the number of minutes per tick
// 0 means no automatic server tick
func OptionTick(minutes int) ServerOption {
	return func(s *Server) {
		s.minutesPerTick = minutes
	}
}

// NewServer sets up a new server
func NewServer(opts ...ServerOption) *Server {

	// Set up the default server
	s := &Server{
		address:     "",
		persistence: EphemeralData,
		schedule:    cron.New(),
		world:       rove.NewWorld(32),
		accountant:  NewSimpleAccountant(),
	}

	// Apply all options
	for _, o := range opts {
		o(s)
	}

	return s
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise(fillWorld bool) (err error) {

	// Add to our sync
	s.sync.Add(1)

	// Load the world file
	if err := s.LoadWorld(); err != nil {
		return err
	}

	// Set up the RPC server and register
	s.netListener, err = net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServ = grpc.NewServer()
	roveapi.RegisterRoveServer(s.grpcServ, s)

	return nil
}

// Addr will return the server address set after the listen
func (s *Server) Addr() string {
	return s.address
}

// Run executes the server
func (s *Server) Run() {
	defer s.sync.Done()

	// Set up the schedule if requested
	if s.minutesPerTick != 0 {
		if err := s.schedule.AddFunc(fmt.Sprintf("0 */%d * * *", s.minutesPerTick), func() {
			// Ensure we don't quit during this function
			s.sync.Add(1)
			defer s.sync.Done()

			log.Println("Executing server tick")

			// Tick the world
			s.world.Tick()

			// Save out the new world state
			if err := s.SaveWorld(); err != nil {
				log.Fatalf("Failed to save the world: %s", err)
			}
		}); err != nil {
			log.Fatal(err)
		}
		s.schedule.Start()
		log.Printf("First server tick scheduled for %s\n", s.schedule.Entries()[0].Next.Format("15:04:05"))
	}

	// Serve the RPC server
	log.Printf("Serving gRPC on %s\n", s.address)
	if err := s.grpcServ.Serve(s.netListener); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("failed to serve gRPC: %s", err)
	}
}

// Stop will stop the current server
func (s *Server) Stop() error {
	// Stop the cron
	s.schedule.Stop()

	// Stop the gRPC
	s.grpcServ.Stop()

	return nil
}

// Close waits until the server is finished and closes up shop
func (s *Server) Close() error {
	// Wait until the world has shut down
	s.sync.Wait()

	// Save and return
	return s.SaveWorld()
}

// StopAndClose waits until the server is finished and closes up shop
func (s *Server) StopAndClose() error {
	// Stop the server
	if err := s.Stop(); err != nil {
		return err
	}

	// Close and return
	return s.Close()
}

// SaveWorld will save out the world file
func (s *Server) SaveWorld() error {
	if s.persistence == PersistentData {
		s.world.RLock()
		defer s.world.RUnlock()
		if err := persistence.SaveAll("world", s.world, "accounts", s.accountant); err != nil {
			return fmt.Errorf("failed to save out persistent data: %s", err)
		}
	}
	return nil
}

// LoadWorld will load all persistent data
func (s *Server) LoadWorld() error {
	if s.persistence == PersistentData {
		s.world.Lock()
		defer s.world.Unlock()
		if err := persistence.LoadAll("world", &s.world, "accounts", &s.accountant); err != nil {
			return err
		}
	}
	return nil
}

// SpawnRoverForAccount spawns the rover rover for an account
func (s *Server) SpawnRoverForAccount(account string) (string, error) {
	inst, err := s.world.SpawnRover()
	if err != nil {
		return "", err
	}

	err = s.accountant.AssignData(account, "rover", inst)
	if err != nil {
		log.Printf("Failed to assign rover to account, %s", err)

		// Try and clear up the rover
		if err := s.world.DestroyRover(inst); err != nil {
			log.Printf("Failed to destroy rover after failed rover assign: %s", err)
		}

		return "", err
	}

	return inst, nil
}
