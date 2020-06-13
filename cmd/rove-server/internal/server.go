package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
)

var accountantAddress = os.Getenv("ACCOUNTANT_ADDRESS")

const (
	// PersistentData will allow the server to load and save it's state
	PersistentData = iota

	// EphemeralData will let the server neither load or save out any of it's data
	EphemeralData
)

// Server contains the relevant data to run a game server
type Server struct {

	// Internal state
	world *game.World

	// Accountant server
	accountant accounts.AccountantClient
	clientConn *grpc.ClientConn

	// gRPC server
	netListener net.Listener
	grpcServ    *grpc.Server

	// Config settings
	address     string
	persistence int
	tick        int

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
		s.tick = minutes
	}
}

// NewServer sets up a new server
func NewServer(opts ...ServerOption) *Server {

	// Set up the default server
	s := &Server{
		address:     "",
		persistence: EphemeralData,
		schedule:    cron.New(),
	}

	// Apply all options
	for _, o := range opts {
		o(s)
	}

	// Start small, we can grow the world later
	s.world = game.NewWorld(4, 8)

	return s
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise(fillWorld bool) (err error) {

	// Add to our sync
	s.sync.Add(1)

	// Connect to the accountant
	if len(accountantAddress) == 0 {
		log.Fatal("must set ACCOUNTANT_ADDRESS")
	}
	log.Printf("Dialing accountant on %s\n", accountantAddress)
	s.clientConn, err = grpc.Dial(accountantAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.accountant = accounts.NewAccountantClient(s.clientConn)

	// Spawn a border on the default world
	if err := s.world.SpawnWorld(fillWorld); err != nil {
		return err
	}

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
	rove.RegisterRoveServer(s.grpcServ, s)

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
	if s.tick != 0 {
		if err := s.schedule.AddFunc(fmt.Sprintf("0 */%d * * *", s.tick), func() {
			// Ensure we don't quit during this function
			s.sync.Add(1)
			defer s.sync.Done()

			log.Println("Executing server tick")

			// Run the command queues
			s.world.ExecuteCommandQueues()

			// Save out the new world state
			s.SaveWorld()
		}); err != nil {
			log.Fatal(err)
		}
		s.schedule.Start()
		log.Printf("First server tick scheduled for %s\n", s.schedule.Entries()[0].Next.Format("15:04:05"))
	}

	// Serve the RPC server
	log.Printf("Serving rove on %s\n", s.address)
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

	// Close the accountant connection
	if err := s.clientConn.Close(); err != nil {
		return err
	}

	return nil
}

// Close waits until the server is finished and closes up shop
func (s *Server) Close() error {
	// Wait until the world has shut down
	s.sync.Wait()

	// Save and return
	return s.SaveWorld()
}

// Close waits until the server is finished and closes up shop
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
		if err := persistence.SaveAll("world", s.world); err != nil {
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
		if err := persistence.LoadAll("world", &s.world); err != nil {
			return err
		}
	}
	return nil
}

// used as the type for the return struct
type BadRequestError struct {
	Error string `json:"error"`
}

// SpawnRoverForAccount spawns the rover rover for an account
func (s *Server) SpawnRoverForAccount(account string) (game.RoverAttributes, uuid.UUID, error) {
	if inst, err := s.world.SpawnRover(); err != nil {
		return game.RoverAttributes{}, uuid.UUID{}, err

	} else if attribs, err := s.world.RoverAttributes(inst); err != nil {
		return game.RoverAttributes{}, uuid.UUID{}, fmt.Errorf("no attributes found for created rover: %s", err)

	} else {
		keyval := accounts.DataKeyValue{Account: account, Key: "rover", Value: inst.String()}
		_, err := s.accountant.AssignValue(context.Background(), &keyval)
		if err != nil {
			log.Printf("Failed to assign rover to account, %s", err)

			// Try and clear up the rover
			if err := s.world.DestroyRover(inst); err != nil {
				log.Printf("Failed to destroy rover after failed rover assign: %s", err)
			}

			return game.RoverAttributes{}, uuid.UUID{}, err
		} else {
			return attribs, inst, nil
		}
	}
}
