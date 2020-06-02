package game

import (
	"fmt"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	instances map[uuid.UUID]Instance
}

// Instance describes a single entity or instance of an entity in the world
type Instance struct {
	// id is a unique ID for this instance
	id uuid.UUID

	// pos represents where this instance is in the world
	pos Position
}

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		instances: make(map[uuid.UUID]Instance),
	}
}

// Adds an instance to the game
func (w *World) CreateInstance() uuid.UUID {
	id := uuid.New()

	// Initialise the instance
	instance := Instance{
		id: id,
	}

	// Append the instance to the list
	w.instances[id] = instance

	return id
}

// GetPosition returns the position of a given instance
func (w World) GetPosition(id uuid.UUID) (Position, error) {
	if i, ok := w.instances[id]; ok {
		return i.pos, nil
	} else {
		return Position{}, fmt.Errorf("no instance matching id")
	}
}
