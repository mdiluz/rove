package game

import (
	"fmt"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	Instances map[uuid.UUID]Instance `json:"instances"`

	// dataPath is the location for the data to be stored
	dataPath string
}

// Instance describes a single entity or instance of an entity in the world
type Instance struct {
	// id is a unique ID for this instance
	id uuid.UUID

	// pos represents where this instance is in the world
	pos Position
}

const kWorldFileName = "rove-world.json"

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		Instances: make(map[uuid.UUID]Instance),
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
	w.Instances[id] = instance

	return id
}

// Removes an instance from the game
func (w *World) DestroyInstance(id uuid.UUID) error {
	if _, ok := w.Instances[id]; ok {
		delete(w.Instances, id)
	} else {
		return fmt.Errorf("no instance matching id")
	}
	return nil
}

// GetPosition returns the position of a given instance
func (w World) GetPosition(id uuid.UUID) (Position, error) {
	if i, ok := w.Instances[id]; ok {
		return i.pos, nil
	} else {
		return Position{}, fmt.Errorf("no instance matching id")
	}
}
