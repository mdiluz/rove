package rove

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
)

const (
	maxLogEntries = 16
)

// RoverLogEntry describes a single log entry for the rover
type RoverLogEntry struct {
	// Time is the timestamp of the entry
	Time time.Time

	// Text contains the information in this log entry
	Text string
}

// Rover describes a single rover in the world
type Rover struct {
	// Unique name of this rover
	Name string

	// Pos represents where this rover is in the world
	Pos maths.Vector

	// Bearing is the current direction the rover is facing
	Bearing roveapi.Bearing

	// Range represents the distance the unit's radar can see
	Range int

	// Inventory represents any items the rover is carrying
	Inventory []Object

	// Capacity is the maximum number of inventory items
	Capacity int

	// Integrity represents current rover health
	Integrity int

	// MaximumIntegrity is the full integrity of the rover
	MaximumIntegrity int

	// Charge is the amount of energy the rover has
	Charge int

	// MaximumCharge is the maximum charge able to be stored
	MaximumCharge int

	// SailPosition is the current position of the sails
	SailPosition roveapi.SailPosition

	// Current number of ticks in this move, used for sailing speeds
	MoveTicks int

	// Logs Stores log of information
	Logs []RoverLogEntry

	// The account that owns this rover
	Owner string
}

// DefaultRover returns a default rover object with default settings
func DefaultRover() *Rover {
	return &Rover{
		Range:            4,
		Integrity:        10,
		MaximumIntegrity: 10,
		Capacity:         10,
		Charge:           10,
		MaximumCharge:    10,
		Bearing:          roveapi.Bearing_North,
		SailPosition:     roveapi.SailPosition_SolarCharging,
		Name:             GenerateRoverName(),
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

	// Limit the number of logs
	if len(r.Logs) > maxLogEntries {
		r.Logs = r.Logs[len(r.Logs)-maxLogEntries:]
	}
}

var wordsFile = os.Getenv("WORDS_FILE")
var roverWords []string

// GenerateRoverName generates a new rover name
func GenerateRoverName() string {

	// Try and load the rover words file
	if len(roverWords) == 0 {
		// Try and load the words file
		if file, err := os.Open(wordsFile); err != nil {
			log.Printf("Couldn't read words file [%s], running without words: %s\n", wordsFile, err)
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				roverWords = append(roverWords, scanner.Text())
			}
			if scanner.Err() != nil {
				log.Printf("Failure during word file scan: %s\n", scanner.Err())
			}
		}
	}

	// Assign a random name if we have words
	if len(roverWords) > 0 {
		// Loop until we find a unique name
		return fmt.Sprintf("%s-%s", roverWords[rand.Intn(len(roverWords))], roverWords[rand.Intn(len(roverWords))])
	}

	// Default to a unique string
	return uuid.New().String()
}
