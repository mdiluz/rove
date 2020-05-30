package rovegame

import "github.com/google/uuid"

// World describes a self contained universe and everything in it
type World struct {
	instances []Instance
}

// Instance describes a single entity or instance of an entity in the world
type Instance struct {
	id uuid.UUID
}

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{}
}

// Adds an instance to the game
func (w *World) CreateInstance() uuid.UUID {
	id := uuid.New()

	// Initialise the instance
	instance := Instance{
		id: id,
	}

	// Append the instance to the list
	w.instances = append(w.instances, instance)

	return id
}
