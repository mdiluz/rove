package internal

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
)

var accountantAddress = flag.String("accountant", "", "address of the accountant to connect to")

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

	// HTTP server
	listener net.Listener
	server   *http.Server
	router   *mux.Router

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

	router := mux.NewRouter().StrictSlash(true)

	// Set up the default server
	s := &Server{
		address:     "",
		persistence: EphemeralData,
		router:      router,
		schedule:    cron.New(),
	}

	// Apply all options
	for _, o := range opts {
		o(s)
	}

	// Set up the server object
	s.server = &http.Server{Addr: s.address, Handler: s.router}

	// Start small, we can grow the world later
	s.world = game.NewWorld(4, 8)

	return s
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise(fillWorld bool) (err error) {

	// Add to our sync
	s.sync.Add(1)

	// Connect to the accountant
	fmt.Printf("Dialing accountant on %s\n", *accountantAddress)
	clientConn, err := grpc.Dial(*accountantAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.accountant = accounts.NewAccountantClient(clientConn)

	// Spawn a border on the default world
	if err := s.world.SpawnWorld(fillWorld); err != nil {
		return err
	}

	// Load the world file
	if err := s.LoadWorld(); err != nil {
		return err
	}

	// Set up the handlers
	for _, route := range Routes {
		s.router.HandleFunc(route.path, s.wrapHandler(route.method, route.handler))
	}

	// Start the listen
	if s.listener, err = net.Listen("tcp", s.server.Addr); err != nil {
		return err
	}

	s.address = s.listener.Addr().String()
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

			fmt.Println("Executing server tick")

			// Run the command queues
			s.world.ExecuteCommandQueues()

			// Save out the new world state
			s.SaveWorld()
		}); err != nil {
			log.Fatal(err)
		}
		s.schedule.Start()
		fmt.Printf("First server tick scheduled for %s\n", s.schedule.Entries()[0].Next.Format("15:04:05"))
	}

	// Serve the http requests
	if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Stop will stop the current server
func (s *Server) Stop() error {
	// Stop the cron
	s.schedule.Stop()

	// Try and shut down the http server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	// Close the accountant connection
	if err := s.clientConn.Close(); err != nil {
		return err
	}

	return nil
}

// Close waits until the server is finished and closes up shop
func (s *Server) Close() error {
	// Wait until the server has shut down
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

// wrapHandler wraps a request handler in http checks
func (s *Server) wrapHandler(method string, handler Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

		vars := mux.Vars(r)

		// Verify the method, call the handler, and encode the return
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)

		} else if val, err := handler(s, vars, r.Body, w); err != nil {
			fmt.Printf("Failed to handle http request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

		} else if err := json.NewEncoder(w).Encode(val); err != nil {
			fmt.Printf("Failed to encode reply to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
	}
}

// SpawnRoverForAccount spawns the rover rover for an account
func (s *Server) SpawnRoverForAccount(account string) (game.RoverAttributes, uuid.UUID, error) {
	if inst, err := s.world.SpawnRover(); err != nil {
		return game.RoverAttributes{}, uuid.UUID{}, err

	} else if attribs, err := s.world.RoverAttributes(inst); err != nil {
		return game.RoverAttributes{}, uuid.UUID{}, fmt.Errorf("No attributes found for created rover: %s", err)

	} else {
		keyval := accounts.DataKeyValue{Account: account, Key: "rover", Value: inst.String()}
		resp, err := s.accountant.AssignValue(context.Background(), &keyval)
		if err != nil || !resp.Success {
			fmt.Printf("Failed to assign rover to account, %s, %s", err, resp.Error)

			// Try and clear up the rover
			if err := s.world.DestroyRover(inst); err != nil {
				fmt.Printf("Failed to destroy rover after failed rover assign: %s", err)
			}

			return game.RoverAttributes{}, uuid.UUID{}, err
		} else {
			return attribs, inst, nil
		}
	}
}
