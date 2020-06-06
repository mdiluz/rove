package game

import "github.com/google/uuid"

// RoverAttributes contains attributes of a rover
type RoverAttributes struct {
	// Speed represents the Speed that the rover will move per second
	Speed int `json:"speed"`

	// Range represents the distance the unit's radar can see
	Range int `json:"range"`

	// Name of this rover
	Name string

	// Pos represents where this rover is in the world
	Pos Vector `json:"pos"`
}

// Rover describes a single rover in the world
type Rover struct {
	// Id is a unique ID for this rover
	Id uuid.UUID `json:"id"`

	// Attributes represents the physical attributes of the rover
	Attributes RoverAttributes `json:"attributes"`
}
