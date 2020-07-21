package rove

import (
	"testing"

	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Toggle(t *testing.T) {
	w := NewWorld(8)
	a, err := w.SpawnRover()
	assert.NoError(t, err)

	r, err := w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_SolarCharging, r.SailPosition)

	w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_toggle})
	w.EnqueueAllIncoming()
	w.ExecuteCommandQueues()

	r, err = w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_CatchingWind, r.SailPosition)

	w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_toggle})
	w.EnqueueAllIncoming()
	w.ExecuteCommandQueues()

	r, err = w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_SolarCharging, r.SailPosition)
}

func TestCommand_Turn(t *testing.T) {
	w := NewWorld(8)
	a, err := w.SpawnRover()
	assert.NoError(t, err)

	w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_turn, Data: &roveapi.Command_Turn{Turn: roveapi.Bearing_NorthWest}})
	w.EnqueueAllIncoming()
	w.ExecuteCommandQueues()

	r, err := w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_NorthWest, r.Bearing)
}

func TestCommand_Stash(t *testing.T) {
	// TODO: Test the stash command
}

func TestCommand_Repair(t *testing.T) {
	// TODO: Test the repair command
}

func TestCommand_Broadcast(t *testing.T) {
	// TODO: Test the stash command
}

func TestCommand_Invalid(t *testing.T) {
	// TODO: Test the invalid command
}
