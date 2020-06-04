package game

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Spawn(t *testing.T) {
	world := NewWorld()
	a := uuid.New()

	spawnCommand := world.CommandSpawn(a)
	assert.NoError(t, world.Execute(spawnCommand), "Failed to execute spawn command")

	rover, ok := world.Rovers[a]
	assert.True(t, ok, "No new rover in world")
	assert.Equal(t, a, rover.Id, "New rover has incorrect id")
}

func TestCommand_Move(t *testing.T) {
	world := NewWorld()
	a := uuid.New()
	assert.NoError(t, world.SpawnRover(a), "Failed to spawn")

	pos := Vector{
		X: 1.0,
		Y: 2.0,
	}

	err := world.SetPosition(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	// TODO: Test the bearing/duration movement
	/*
		// Try the move command
		moveCommand := world.CommandMove(a, move)
		assert.NoError(t, world.Execute(moveCommand), "Failed to execute move command")

		newpos, err := world.GetPosition(a)
		assert.NoError(t, err, "Failed to set position for rover")
		pos.Add(move)
		assert.Equal(t, pos, newpos, "Failed to correctly set position for rover")
	*/
}
