package game

import "github.com/mdiluz/rove/pkg/rove"

// Command represends a single command to execute
type Command struct {
	Command rove.CommandType `json:"command"`

	// Used in the move command
	Bearing string `json:"bearing,omitempty"`

	// Used in the broadcast command
	Message []byte `json:"message,omitempty"`
}

// CommandStream is a list of commands to execute in order
type CommandStream []Command
