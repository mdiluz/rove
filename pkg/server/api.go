package server

import "github.com/mdiluz/rove/pkg/game"

// ==============================
// API: /status method: GET
// Queries the status of the server

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready   bool   `json:"ready"`
	Version string `json:"version"`
}

// ==============================
// API: /register method: POST
// Registers a user account by name
// Responds with a unique ID for that account to be used in future requests

// RegisterData describes the data to send when registering
type RegisterData struct {
	Name string `json:"name"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`

	Id string `json:"id"`
}

// ==============================
// API: /spawn method: POST
// Spawns the primary entity for an account
// Responds with the position of said entity

// SpawnData is the data to be sent for the spawn command
type SpawnData struct {
	Id string `json:"id"`
}

// SpawnResponse is the data to respond with on a spawn command
type SpawnResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`

	Position game.Vector `json:"position"`
}

// ==============================
// API: /commands method: POST
// Issues a set of commands from the user

// CommandsData is a set of commands to execute in order
type CommandsData struct {
	Id       string    `json:"id"`
	Commands []Command `json:"commands"`
}

// CommandsResponse is the response to be sent back
type CommandsResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
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
	Vector game.Vector `json:"vector"`
}

// ================
// API: /view POST
// Queries the current view for the user

// ViewData describes the input data to request an accounts current view
type ViewData struct {
	Id string `json:"id"`
}

// ViewResponse describes the response to a /view call
type ViewResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
