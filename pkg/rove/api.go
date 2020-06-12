package rove

import (
	"path"

	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/game"
)

// ==============================
// API: /status method: GET

// Status queries the status of the server
func (s Server) Status() (r StatusResponse, err error) {
	s.Get("status", &r)
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
func (s Server) Register(d RegisterData) (r RegisterResponse, err error) {
	err = s.Post("register", d, &r)
	return
}

// RegisterData describes the data to send when registering
type RegisterData struct {
	Name string `json:"name"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	// Placeholder for future information
}

// ==============================
// API: /{account}/command method: POST

// Command issues a set of commands from the user
func (s Server) Command(account string, d CommandData) (r CommandResponse, err error) {
	err = s.Post(path.Join(account, "command"), d, &r)
	return
}

// CommandData is a set of commands to execute in order
type CommandData struct {
	Commands []game.Command `json:"commands"`
}

// CommandResponse is the response to be sent back
type CommandResponse struct {
	// Placeholder for future information
}

// ================
// API: /{account}/radar method: GET

// Radar queries the current radar for the user
func (s Server) Radar(account string) (r RadarResponse, err error) {
	err = s.Get(path.Join(account, "radar"), &r)
	return
}

// RadarResponse describes the response to a /radar call
type RadarResponse struct {
	// The set of positions for nearby non-empty tiles
	Range int          `json:"range"`
	Tiles []atlas.Tile `json:"tiles"`
}

// ================
// API: /{account}/rover method: GET

// Rover queries the current state of the rover
func (s Server) Rover(account string) (r RoverResponse, err error) {
	err = s.Get(path.Join(account, "rover"), &r)
	return
}

// RoverResponse includes information about the rover in question
type RoverResponse struct {
	// The current position of this rover
	Attributes game.RoverAttributes `json:"attributes"`
}
