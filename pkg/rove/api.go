package rove

import (
	"path"

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
	Ready    bool   `json:"ready"`
	Version  string `json:"version"`
	Tick     int    `json:"tick"`
	NextTick string `json:"nexttick,omitempty"`
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
// API: /{account}/spawn method: POST

// Spawn spawns the rover for an account
// Responds with the position of said rover
func (s Server) Spawn(account string, d SpawnData) (r SpawnResponse, err error) {
	err = s.POST(path.Join(account, "spawn"), d, &r)
	return
}

// SpawnData is the data to be sent for the spawn command
type SpawnData struct {
	// Empty for now, reserved for data
}

// SpawnResponse is the data to respond with on a spawn command
type SpawnResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The attributes of the spawned entity
	Attributes game.RoverAttributes `json:"attributes"`
}

// ==============================
// API: /{account}/command method: POST

// Command issues a set of commands from the user
func (s Server) Command(account string, d CommandData) (r CommandResponse, err error) {
	err = s.POST(path.Join(account, "command"), d, &r)
	return
}

// CommandData is a set of commands to execute in order
type CommandData struct {
	Commands []game.Command `json:"commands"`
}

// CommandResponse is the response to be sent back
type CommandResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ================
// API: /{account}/radar method: GET

// Radar queries the current radar for the user
func (s Server) Radar(account string) (r RadarResponse, err error) {
	err = s.GET(path.Join(account, "radar"), &r)
	return
}

// RadarResponse describes the response to a /radar call
type RadarResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The set of positions for nearby rovers
	Rovers []game.Vector `json:"rovers"`
}

// ================
// API: /{account}/rover method: GET

// Rover queries the current state of the rover
func (s Server) Rover(account string) (r RoverResponse, err error) {
	err = s.GET(path.Join(account, "rover"), &r)
	return
}

// RoverResponse includes information about the rover in question
type RoverResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// The current position of this rover
	Attributes game.RoverAttributes `json:"attributes"`
}
