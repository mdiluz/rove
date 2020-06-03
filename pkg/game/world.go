package game

import (
	"fmt"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	// Instances is a map of all the instances in the game
	Instances map[uuid.UUID]Instance `json:"instances"`
}

// Instance describes a single entity or instance of an entity in the world
type Instance struct {
	// id is a unique ID for this instance
	id uuid.UUID

	// pos represents where this instance is in the world
	pos Vector
}

const kWorldFileName = "rove-world.json"

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		Instances: make(map[uuid.UUID]Instance),
	}
}

// Spawn adds an instance to the game
func (w *World) Spawn(id uuid.UUID) error {
	if _, ok := w.Instances[id]; ok {
		return fmt.Errorf("instance with id %s already exists in world", id)
	}

	// Initialise the instance
	instance := Instance{
		id: id,
	}

	// Append the instance to the list
	w.Instances[id] = instance

	return nil
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
func (w World) GetPosition(id uuid.UUID) (Vector, error) {
	if i, ok := w.Instances[id]; ok {
		return i.pos, nil
	} else {
		return Vector{}, fmt.Errorf("no instance matching id")
	}
}

// SetPosition sets an instances position
func (w *World) SetPosition(id uuid.UUID, pos Vector) error {
	if i, ok := w.Instances[id]; ok {
		i.pos = pos
		w.Instances[id] = i
		return nil
	} else {
		return fmt.Errorf("no instance matching id")
	}
}

// SetPosition sets an instances position
func (w *World) MovePosition(id uuid.UUID, vec Vector) (Vector, error) {
	if i, ok := w.Instances[id]; ok {
		i.pos.Add(vec)
		w.Instances[id] = i
		return i.pos, nil
	} else {
		return Vector{}, fmt.Errorf("no instance matching id")
	}
}

// Execute will run the commands given
func (w *World) Execute(commands ...Command) error {
	for _, c := range commands {
		if err := c(); err != nil {
			return err
		}
	}
	return nil
}
