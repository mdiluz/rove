package rove

import (
	"testing"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Move(t *testing.T) {
	world := NewWorld(8)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	pos := maths.Vector{
		X: 1.0,
		Y: 2.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	// Try the move command
	moveCommand := Command{Command: roveapi.CommandType_move, Bearing: roveapi.Bearing_North}
	assert.NoError(t, world.Enqueue(a, moveCommand), "Failed to execute move command")

	// Tick the world
	world.EnqueueAllIncoming()
	world.ExecuteCommandQueues()

	newPos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(maths.Vector{X: 0.0, Y: 1})
	assert.Equal(t, pos, newPos, "Failed to correctly set position for rover")
}

func TestCommand_Recharge(t *testing.T) {
	world := NewWorld(8)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	pos := maths.Vector{
		X: 1.0,
		Y: 2.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	// Move to use up some charge
	moveCommand := Command{Command: roveapi.CommandType_move, Bearing: roveapi.Bearing_North}
	assert.NoError(t, world.Enqueue(a, moveCommand), "Failed to queue move command")

	// Tick the world
	world.EnqueueAllIncoming()
	world.ExecuteCommandQueues()

	rover, _ := world.GetRover(a)
	assert.Equal(t, rover.MaximumCharge-1, rover.Charge)

	chargeCommand := Command{Command: roveapi.CommandType_recharge}
	assert.NoError(t, world.Enqueue(a, chargeCommand), "Failed to queue recharge command")

	// Tick the world
	world.EnqueueAllIncoming()
	world.ExecuteCommandQueues()

	rover, _ = world.GetRover(a)
	assert.Equal(t, rover.MaximumCharge, rover.Charge)

}
