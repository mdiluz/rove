package game

import (
	"fmt"
	"math"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	// Rovers is a id->data map of all the rovers in the game
	Rovers map[uuid.UUID]Rover `json:"rovers"`
}

// RoverAttributes contains attributes of a rover
type RoverAttributes struct {
	// Speed represents the Speed that the rover will move per second
	Speed float64 `json:"speed"`

	// Sight represents the distance the unit can see
	Sight float64 `json:"sight"`
}

// Rover describes a single rover in the world
type Rover struct {
	// Id is a unique ID for this rover
	Id uuid.UUID `json:"id"`

	// Pos represents where this rover is in the world
	Pos Vector `json:"pos"`

	// Attributes represents the physical attributes of the rover
	Attributes RoverAttributes `json:"attributes"`
}

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		Rovers: make(map[uuid.UUID]Rover),
	}
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover() uuid.UUID {
	// Initialise the rover
	rover := Rover{
		Id: uuid.New(),

		// TODO: Set this somehow
		Pos: Vector{},

		// TODO: Stop these being random numbers
		Attributes: RoverAttributes{
			Speed: 1.0,
			Sight: 20.0,
		},
	}

	// Append the rover to the list
	w.Rovers[rover.Id] = rover

	return rover.Id
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

// RoverAttributes returns the attributes of a requested rover
func (w World) RoverAttributes(id uuid.UUID) (RoverAttributes, error) {
	if i, ok := w.Rovers[id]; ok {
		return i.Attributes, nil
	} else {
		return RoverAttributes{}, fmt.Errorf("no rover matching id")
	}
}

// RoverPosition returns the position of a given rover
func (w World) RoverPosition(id uuid.UUID) (Vector, error) {
	if i, ok := w.Rovers[id]; ok {
		return i.Pos, nil
	} else {
		return Vector{}, fmt.Errorf("no rover matching id")
	}
}

// WarpRover sets an rovers position
func (w *World) WarpRover(id uuid.UUID, pos Vector) error {
	if i, ok := w.Rovers[id]; ok {
		i.Pos = pos
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) MoveRover(id uuid.UUID, bearing float64, duration float64) (Vector, error) {
	if i, ok := w.Rovers[id]; ok {
		// Calculate the distance
		distance := i.Attributes.Speed * float64(duration)

		// Calculate the full movement based on the bearing
		move := Vector{
			X: math.Sin(bearing) * distance,
			Y: math.Cos(bearing) * distance,
		}

		// Increment the position by the movement
		i.Pos.Add(move)

		// Set the rover values to the new ones
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
