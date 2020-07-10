package rove

import "github.com/mdiluz/rove/pkg/roveapi"

// Command represends a single command to execute
type Command struct {
	Command roveapi.CommandType `json:"command"`

	// Used in the move command
	Bearing string `json:"bearing,omitempty"`

	// Used in the broadcast command
	Message []byte `json:"message,omitempty"`
}

// CommandStream is a list of commands to execute in order
type CommandStream []Command
