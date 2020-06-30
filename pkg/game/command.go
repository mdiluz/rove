package game

const (
	// CommandMove Moves the rover in the chosen bearing
	CommandMove = "move"

	// CommandStash Will attempt to stash the object at the current location
	CommandStash = "stash"

	// CommandRepair Will attempt to repair the rover with an inventory object
	CommandRepair = "repair"
)

// Command represends a single command to execute
type Command struct {
	Command string `json:"command"`

	// Used in the move command
	Bearing string `json:"bearing,omitempty"`
}

// CommandStream is a list of commands to execute in order
type CommandStream []Command
