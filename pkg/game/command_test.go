package game

import (
	"testing"

	"github.com/mdiluz/rove/pkg/bearing"
	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Move(t *testing.T) {
	world := NewWorld(2, 8)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	pos := vector.Vector{
		X: 1.0,
		Y: 2.0,
	}

	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	bearing := bearing.North
	duration := 1
	// Try the move command
	moveCommand := Command{Command: CommandMove, Bearing: bearing.String(), Duration: duration}
	assert.NoError(t, world.Enqueue(a, moveCommand), "Failed to execute move command")

	// Tick the world
	world.EnqueueAllIncoming()
	world.ExecuteCommandQueues()

	newatributes, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(vector.Vector{X: 0.0, Y: duration * attribs.Speed}) // We should have moved duration*speed north
	assert.Equal(t, pos, newatributes.Pos, "Failed to correctly set position for rover")
}
