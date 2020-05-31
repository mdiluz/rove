package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/game"
)

// Server contains the relevant data to run a game server
type Server struct {
	port int

	accountant *Accountant
	world      *game.World

	router *mux.Router
}

// NewServer sets up a new server
func NewServer(port int) *Server {
	return &Server{
		port:       port,
		accountant: NewAccountant(),
		world:      game.NewWorld(),
	}
}

// Initialise sets up internal state ready to serve
func (s *Server) Initialise() {
	// Set up the world
	s.world = game.NewWorld()
	fmt.Printf("World created\n\t%+v\n", s.world)

	// Create a new router
	s.SetUpRouter()
	fmt.Printf("Routes Created\n")
}

// Run executes the server
func (s *Server) Run() {
	// Listen and serve the http requests
	fmt.Println("Serving HTTP")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router); err != nil {
		log.Fatal(err)
	}
}
