package rove

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/maths"
)

// RoverLogEntry describes a single log entry for the rover
type RoverLogEntry struct {
	// Time is the timestamp of the entry
	Time time.Time `json:"time"`

	// Text contains the information in this log entry
	Text string `json:"text"`
}

// Rover describes a single rover in the world
type Rover struct {
	// Unique name of this rover
	Name string `json:"name"`

	// Pos represents where this rover is in the world
	Pos maths.Vector `json:"pos"`

	// Range represents the distance the unit's radar can see
	Range int `json:"range"`

	// Inventory represents any items the rover is carrying
	Inventory []Object `json:"inventory"`

	// Capacity is the maximum number of inventory items
	Capacity int `json:"capacity"`

	// Integrity represents current rover health
	Integrity int `json:"integrity"`

	// MaximumIntegrity is the full integrity of the rover
	MaximumIntegrity int `json:"maximum-integrity"`

	// Charge is the amount of energy the rover has
	Charge int `json:"charge"`

	// MaximumCharge is the maximum charge able to be stored
	MaximumCharge int `json:"maximum-Charge"`

	// Logs Stores log of information
	Logs []RoverLogEntry `json:"logs"`
}

// DefaultRover returns a default rover object with default settings
func DefaultRover() Rover {
	return Rover{
		Range:            4,
		Integrity:        10,
		MaximumIntegrity: 10,
		Capacity:         10,
		Charge:           10,
		MaximumCharge:    10,
		Name:             uuid.New().String(),
	}
}

// AddLogEntryf adds an entry to the rovers log
func (r *Rover) AddLogEntryf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	log.Printf("%s log entry: %s", r.Name, text)
	r.Logs = append(r.Logs,
		RoverLogEntry{
			Time: time.Now(),
			Text: text,
		},
	)
}
