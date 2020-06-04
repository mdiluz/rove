package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Move(t *testing.T) {
	world := NewWorld()
	a := world.SpawnRover()
	pos := Vector{
		X: 1.0,
		Y: 2.0,
	}

	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	bearing := 0.0
	duration := 1.0
	// Try the move command
	moveCommand := world.CommandMove(a, bearing, duration)
	assert.NoError(t, world.Execute(moveCommand), "Failed to execute move command")

	newpos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(Vector{0.0, float64(duration) * attribs.Speed}) // We should have moved duration*speed north
	assert.Equal(t, pos, newpos, "Failed to correctly set position for rover")
}
