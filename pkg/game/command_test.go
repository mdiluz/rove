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

	bearing := North
	duration := 1
	// Try the move command
	moveCommand := Command{Command: CommandMove, Bearing: bearing.String(), Duration: duration}
	assert.NoError(t, world.Execute(a, moveCommand), "Failed to execute move command")

	newpos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(Vector{0.0, duration * attribs.Speed}) // We should have moved duration*speed north
	assert.Equal(t, pos, newpos, "Failed to correctly set position for rover")
}
