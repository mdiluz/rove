package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/game"
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

	persistence         int
	persistenceLocation string

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
func OptionPersistentData(loc string) ServerOption {
	return func(s *Server) {
		s.persistence = PersistentData
		s.persistenceLocation = loc
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
	s.accountant = accounts.NewAccountant(s.persistenceLocation)
	s.world = game.NewWorld(s.persistenceLocation)

	return s
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise() error {

	// Load the accounts if requested
	if s.persistence == PersistentData {
		if err := s.accountant.Load(); err != nil {
			return err
		}
		if err := s.world.Load(); err != nil {
			return err
		}
	}

	// Create a new router
	s.SetUpRouter()
	fmt.Printf("Routes Created\n")

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
		if err := s.accountant.Save(); err != nil {
			return err
		}
		if err := s.world.Save(); err != nil {
			return err
		}
	}
	return nil
}
