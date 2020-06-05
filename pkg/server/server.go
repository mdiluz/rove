package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/rove"
)

const (
	// PersistentData will allow the server to load and save it's state
	PersistentData = iota

	// EphemeralData will let the server neither load or save out any of it's data
	EphemeralData
)

// Server contains the relevant data to run a game server
type Server struct {
	port int

	accountant *accounts.Accountant
	world      *game.World

	server *http.Server
	router *mux.Router

	persistence int

	sync sync.WaitGroup
}

// ServerOption defines a server creation option
type ServerOption func(s *Server)

// OptionPort sets the server port for hosting
func OptionPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

// OptionPersistentData sets the server data to be persistent
func OptionPersistentData() ServerOption {
	return func(s *Server) {
		s.persistence = PersistentData
	}
}

// NewServer sets up a new server
func NewServer(opts ...ServerOption) *Server {

	router := mux.NewRouter().StrictSlash(true)

	// Set up the default server
	s := &Server{
		port:        8080,
		persistence: EphemeralData,
		router:      router,
	}

	// Apply all options
	for _, o := range opts {
		o(s)
	}

	// Set up the server object
	s.server = &http.Server{Addr: fmt.Sprintf(":%d", s.port), Handler: router}

	// Create the accountant
	s.accountant = accounts.NewAccountant()
	s.world = game.NewWorld()

	return s
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise() error {

	// Load the accounts if requested
	if s.persistence == PersistentData {
		if err := persistence.LoadAll("accounts", &s.accountant, "world", &s.world); err != nil {
			return err
		}
	}

	// Set up the handlers
	for _, route := range Routes {
		s.router.HandleFunc(route.path, s.wrapHandler(route.method, route.handler))
	}

	// Add to our sync
	s.sync.Add(1)

	return nil
}

// Run executes the server
func (s *Server) Run() {
	defer s.sync.Done()

	// Listen and serve the http requests
	fmt.Printf("Serving HTTP on port %d\n", s.port)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Close closes up the server
func (s *Server) Close() error {
	// Try and shut down the http server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	// Wait until the server is shut down
	s.sync.Wait()

	// Save the accounts if requested
	if s.persistence == PersistentData {
		if err := persistence.SaveAll("accounts", s.accountant, "world", s.world); err != nil {
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

		// Verify the method, call the handler, and encode the return
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)

		} else if val, err := handler(s, r.Body, w); err != nil {
			fmt.Printf("Failed to handle http request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

		} else if err := json.NewEncoder(w).Encode(val); err != nil {
			fmt.Printf("Failed to encode return to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
	}
}

// SpawnRoverForAccount spawns the rover rover for an account
func (s *Server) SpawnRoverForAccount(accountid uuid.UUID) (game.Vector, uuid.UUID, error) {
	inst := s.world.SpawnRover()
	if pos, err := s.world.RoverPosition(inst); err != nil {
		return game.Vector{}, uuid.UUID{}, fmt.Errorf("No position found for created rover")

	} else {
		if err := s.accountant.AssignRover(accountid, inst); err != nil {
			// Try and clear up the rover
			if err := s.world.DestroyRover(inst); err != nil {
				fmt.Printf("Failed to destroy rover after failed rover assign: %s", err)
			}

			return game.Vector{}, uuid.UUID{}, err
		} else {
			return pos, inst, nil
		}
	}
}

// ConvertCommands converts server commands to game commands
func (s *Server) ConvertCommands(commands []rove.Command, inst uuid.UUID) ([]game.Command, error) {
	var cmds []game.Command
	for _, c := range commands {
		switch c.Command {
		case rove.CommandMove:
			if bearing, err := game.DirectionFromString(c.Bearing); err != nil {
				return nil, err
			} else {
				cmds = append(cmds, s.world.CommandMove(inst, bearing, c.Duration))
			}
		}
	}
	return cmds, nil
}
