package game

import "github.com/google/uuid"

// A command is simply a function that acts on the a given rover in the world
type Command func() error

// CommandMove will move the rover in question
func (w *World) CommandMove(id uuid.UUID, bearing float64, duration int64) Command {
	return func() error {
		// TODO: Calculate the move itself

		//_, err := w.MovePosition(id, vec)
		return nil
	}
}

// CommandSpawn
// TODO: Two spawn commands with the same id could trigger a fail later on, we should prevent that somehow
func (w *World) CommandSpawn(id uuid.UUID) Command {
	return func() error {
		return w.SpawnRover(id)
	}
}
