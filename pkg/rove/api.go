package rove

import (
	"github.com/mdiluz/rove/pkg/game"
)

// ==============================
// API: /status method: GET

// Status queries the status of the server
func (s Server) Status() (r StatusResponse, err error) {
	s.GET("status", &r)
	return
}

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready   bool   `json:"ready"`
	Version string `json:"version"`
}

// ==============================
// API: /register method: POST

// Register registers a user by name
// Responds with a unique ID for that user to be used in future requests
func (s Server) Register(d RegisterData) (r RegisterResponse, err error) {
	err = s.POST("register", d, &r)
	return
}

// RegisterData describes the data to send when registering
type RegisterData struct {
	Name string `json:"name"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	Id string `json:"id"`
}

// ==============================
// API: /spawn method: POST

// Spawn spawns the rover for an account
// Responds with the position of said rover
func (s Server) Spawn(d SpawnData) (r SpawnResponse, err error) {
	err = s.POST("spawn", d, &r)
	return
}

// SpawnData is the data to be sent for the spawn command
type SpawnData struct {
	Id string `json:"id"`
}

// SpawnResponse is the data to respond with on a spawn command
type SpawnResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The location of the spawned entity
	Position game.Vector `json:"position"`
}

// ==============================
// API: /commands method: POST

// Commands issues a set of commands from the user
func (s Server) Commands(d CommandsData) (r CommandsResponse, err error) {
	err = s.POST("commands", d, &r)
	return
}

// CommandsData is a set of commands to execute in order
type CommandsData struct {
	Id       string    `json:"id"`
	Commands []Command `json:"commands"`
}

// CommandsResponse is the response to be sent back
type CommandsResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

const (
	// CommandMove describes a single move command
	CommandMove = "move"
)

// Command describes a single command to execute
// it contains the type, and then any members used for each command type
type Command struct {
	// Command is the main command string
	Command string `json:"command"`

	// Used for CommandMove
	Bearing  string `json:"bearing"`  // The direction to move on a compass in short (NW) or long (NorthWest) form
	Duration int    `json:"duration"` // The duration of the move in ticks
}

// ================
// API: /radar POST

// Radar queries the current radar for the user
func (s Server) Radar(d RadarData) (r RadarResponse, err error) {
	err = s.POST("radar", d, &r)
	return
}

// RadarData describes the input data to request an accounts current radar
type RadarData struct {
	Id string `json:"id"`
}

// RadarResponse describes the response to a /radar call
type RadarResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The set of positions for nearby rovers
	Rovers []game.Vector `json:"rovers"`
}

// ================
// API: /rover POST

// Rover queries the current state of the rover
func (s Server) Rover(d RoverData) (r RoverResponse, err error) {
	err = s.POST("rover", d, &r)
	return
}

// RoverData describes the input data to request rover status
type RoverData struct {
	Id string `json:"id"`
}

// RoverResponse includes information about the rover in question
type RoverResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The set of positions for nearby rovers
	Position game.Vector `json:"position"`
}
