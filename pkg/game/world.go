package game

import (
	"fmt"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	// Rovers is a id->data map of all the rovers in the game
	Rovers map[uuid.UUID]Rover `json:"rovers"`
}

// Rover describes a single rover in the world
type Rover struct {
	// Id is a unique ID for this rover
	Id uuid.UUID `json:"id"`

	// Pos represents where this rover is in the world
	Pos Vector `json:"pos"`

	// Speed represents the Speed that the rover will move per second
	Speed float64 `json:"speed"`

	// Sight represents the distance the unit can see
	Sight float64 `json:"sight"`
}

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		Rovers: make(map[uuid.UUID]Rover),
	}
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover(id uuid.UUID) error {
	if _, ok := w.Rovers[id]; ok {
		return fmt.Errorf("rover with id %s already exists in world", id)
	}

	// Initialise the rover
	rover := Rover{
		Id: id,
	}

	// Append the rover to the list
	w.Rovers[id] = rover

	return nil
}

// Removes an rover from the game
func (w *World) DestroyRover(id uuid.UUID) error {
	if _, ok := w.Rovers[id]; ok {
		delete(w.Rovers, id)
	} else {
		return fmt.Errorf("no rover matching id")
	}
	return nil
}

// GetPosition returns the position of a given rover
func (w World) GetPosition(id uuid.UUID) (Vector, error) {
	if i, ok := w.Rovers[id]; ok {
		return i.Pos, nil
	} else {
		return Vector{}, fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) SetPosition(id uuid.UUID, pos Vector) error {
	if i, ok := w.Rovers[id]; ok {
		i.Pos = pos
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) MovePosition(id uuid.UUID, vec Vector) (Vector, error) {
	if i, ok := w.Rovers[id]; ok {
		i.Pos.Add(vec)
		w.Rovers[id] = i
		return i.Pos, nil
	} else {
		return Vector{}, fmt.Errorf("no rover matching id")
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
