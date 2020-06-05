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

// RoverAttributes contains attributes of a rover
type RoverAttributes struct {
	// Speed represents the Speed that the rover will move per second
	Speed int `json:"speed"`

	// Range represents the distance the unit's radar can see
	Range int `json:"range"`
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
			Range: 20.0,
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
func (w *World) MoveRover(id uuid.UUID, bearing Direction, duration int) (Vector, error) {
	if i, ok := w.Rovers[id]; ok {
		// Calculate the distance
		distance := i.Attributes.Speed * duration

		// Calculate the full movement based on the bearing
		move := bearing.Vector().Multiplied(distance)

		// Increment the position by the movement
		i.Pos.Add(move)

		// Set the rover values to the new ones
		w.Rovers[id] = i
		return i.Pos, nil
	} else {
		return Vector{}, fmt.Errorf("no rover matching id")
	}
}

// RadarDescription describes what a rover can see
type RadarDescription struct {
	// Rovers is the set of rovers that this radar can see
	Rovers []Vector `json:"rovers"`
}

// RadarFromRover can be used to query what a rover can currently see
func (w World) RadarFromRover(id uuid.UUID) (RadarDescription, error) {
	if r1, ok := w.Rovers[id]; ok {
		var desc RadarDescription

		// Gather nearby rovers within the range
		for _, r2 := range w.Rovers {
			if r1.Id != r2.Id && r1.Pos.Distance(r2.Pos) < float64(r1.Attributes.Range) {
				desc.Rovers = append(desc.Rovers, r2.Pos)
			}
		}

		return desc, nil
	} else {
		return RadarDescription{}, fmt.Errorf("no rover matching id")
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
