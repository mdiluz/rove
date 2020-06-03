package game

import "github.com/google/uuid"

// A command is simply a function that acts on the a given instance in the world
type Command func() error

// CommandMove will move the instance in question
func (w *World) CommandMove(id uuid.UUID, vec Vector) Command {
	return func() error {
		// Move the instance
		_, err := w.MovePosition(id, vec)
		return err
	}
}

// CommandSpawn
// TODO: Two spawn commands with the same id could trigger a fail later on, we should prevent that somehow
func (w *World) CommandSpawn(id uuid.UUID) Command {
	return func() error {
		return w.Spawn(id)
	}
}
