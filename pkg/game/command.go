package game

import "github.com/google/uuid"

// A command is simply a function that acts on the a given rover in the world
type Command func() error

// CommandMove will move the rover in question
func (w *World) CommandMove(id uuid.UUID, bearing float64, duration float64) Command {
	return func() error {
		_, err := w.MoveRover(id, bearing, duration)
		return err
	}
}
