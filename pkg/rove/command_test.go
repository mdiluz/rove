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

	err = w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_toggle})
	assert.NoError(t, err)
	w.Tick()

	r, err = w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_CatchingWind, r.SailPosition)

	err = w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_toggle})
	assert.NoError(t, err)
	w.Tick()

	r, err = w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_SolarCharging, r.SailPosition)
}

func TestCommand_Turn(t *testing.T) {
	w := NewWorld(8)
	a, err := w.SpawnRover()
	assert.NoError(t, err)

	err = w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_turn, Turn: roveapi.Bearing_NorthWest})
	assert.NoError(t, err)
	w.Tick()

	r, err := w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_NorthWest, r.Bearing)
}

func TestCommand_Stash(t *testing.T) {
	w := NewWorld(8)
	name, err := w.SpawnRover()
	assert.NoError(t, err)

	info, err := w.GetRover(name)
	assert.NoError(t, err)
	assert.Empty(t, info.Inventory)

	// Drop a pickup below us
	w.Atlas.SetObject(info.Pos, Object{Type: roveapi.Object_RockSmall})

	// Try and stash it
	err = w.Enqueue(name, &roveapi.Command{Command: roveapi.CommandType_stash})
	assert.NoError(t, err)
	w.Tick()

	// Check we now have it in the inventory
	info, err = w.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(info.Inventory))
	assert.Equal(t, Object{Type: roveapi.Object_RockSmall}, info.Inventory[0])

	// Check it's no longer on the atlas
	_, obj := w.Atlas.QueryPosition(info.Pos)
	assert.Equal(t, Object{Type: roveapi.Object_ObjectUnknown}, obj)
}

func TestCommand_Repair(t *testing.T) {
	// TODO: Test the repair command
}

func TestCommand_Broadcast(t *testing.T) {
	// TODO: Test the stash command
}

func TestCommand_Invalid(t *testing.T) {
	// TODO: Test an invalid command
}
