package game

import (
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
)

// Rover describes a single rover in the world
type Rover struct {
	// Unique name of this rover
	Name string `json:"name"`

	// Pos represents where this rover is in the world
	Pos vector.Vector `json:"pos"`

	// Range represents the distance the unit's radar can see
	Range int `json:"range"`

	// Inventory represents any items the rover is carrying
	Inventory []objects.Object `json:"inventory"`

	// Capacity is the maximum number of inventory items
	Capacity int `json:"capacity"`

	// Integrity represents current rover health
	Integrity int `json:"integrity"`
}
