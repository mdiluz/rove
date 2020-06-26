package game

import (
	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/vector"
)

// RoverAttributes contains attributes of a rover
type RoverAttributes struct {
	// Name of this rover
	Name string `json:"name"`

	// Range represents the distance the unit's radar can see
	Range int `json:"range"`
}

// Rover describes a single rover in the world
type Rover struct {
	// Id is a unique ID for this rover
	Id uuid.UUID `json:"id"`

	// Pos represents where this rover is in the world
	Pos vector.Vector `json:"pos"`

	// Attributes represents the physical attributes of the rover
	Attributes RoverAttributes `json:"attributes"`

	// Inventory represents any items the rover is carrying
	Inventory []Item `json:"inventory"`
}
