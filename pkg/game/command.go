package game

const (
	CommandMove = "move"
)

// Command represends a single command to execute
type Command struct {
	Command string `json:"command"`

	// Used in the move command
	Bearing  string `json:"bearing,omitempty"`
	Duration int    `json:"duration,omitempty"`
}
