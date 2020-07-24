package rove

import (
	"encoding/json"
	"testing"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Toggle(t *testing.T) {
	w := NewWorld(8)
	a, err := w.SpawnRover("tmp")
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
	a, err := w.SpawnRover("tmp")
	assert.NoError(t, err)

	err = w.Enqueue(a, &roveapi.Command{Command: roveapi.CommandType_turn, Bearing: roveapi.Bearing_NorthWest})
	assert.NoError(t, err)
	w.Tick()

	r, err := w.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_NorthWest, r.Bearing)
}

func TestCommand_Stash(t *testing.T) {
	w := NewWorld(8)
	name, err := w.SpawnRover("tmp")
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
	w := NewWorld(8)
	name, err := w.SpawnRover("tmp")
	assert.NoError(t, err)

	info, err := w.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, info.MaximumIntegrity, info.Integrity)

	// Put a blocking rock to the north
	w.Atlas.SetObject(info.Pos.Added(maths.Vector{X: 0, Y: 1}), Object{Type: roveapi.Object_RockLarge})

	// Try and move and make sure we're blocked
	newpos, err := w.TryMoveRover(name, roveapi.Bearing_North)
	assert.NoError(t, err)
	assert.Equal(t, info.Pos, newpos)

	// Check we're damaged
	info, err = w.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, info.MaximumIntegrity-1, info.Integrity)

	// Stash a repair object
	w.Atlas.SetObject(info.Pos, Object{Type: roveapi.Object_RoverParts})
	obj, err := w.RoverStash(name)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Object_RoverParts, obj)

	// Enqueue the repair and tick
	err = w.Enqueue(name, &roveapi.Command{Command: roveapi.CommandType_repair})
	assert.NoError(t, err)
	w.Tick()

	// Check we're repaired
	info, err = w.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, info.MaximumIntegrity, info.Integrity)
	assert.Equal(t, 0, len(info.Inventory))
}

func TestCommand_Broadcast(t *testing.T) {
	w := NewWorld(8)
	name, err := w.SpawnRover("tmp")
	assert.NoError(t, err)

	// Enqueue the broadcast and tick
	err = w.Enqueue(name, &roveapi.Command{Command: roveapi.CommandType_broadcast, Data: []byte("ABC")})
	assert.NoError(t, err)
	w.Tick()

	info, err := w.GetRover(name)
	assert.NoError(t, err)
	assert.Contains(t, info.Logs[len(info.Logs)-1].Text, "ABC")
}

func TestCommand_Salvage(t *testing.T) {
	w := NewWorld(8)
	name, err := w.SpawnRover("tmp")
	assert.NoError(t, err)

	info, err := w.GetRover(name)
	assert.NoError(t, err)

	w.Atlas.SetObject(info.Pos, Object{Type: roveapi.Object_RoverDormant})

	// Enqueue the broadcast and tick
	err = w.Enqueue(name, &roveapi.Command{Command: roveapi.CommandType_salvage})
	assert.NoError(t, err)
	w.Tick()

	// Check we now have some rover parts
	info, err = w.GetRover(name)
	assert.NoError(t, err)
	assert.NotEmpty(t, info.Inventory)
	for _, i := range info.Inventory {
		assert.Equal(t, roveapi.Object_RoverParts, i.Type)
	}

	// Check the dormant rover is gone
	_, obj := w.Atlas.QueryPosition(info.Pos)
	assert.Equal(t, roveapi.Object_ObjectUnknown, obj.Type)
}

func TestCommand_Transfer(t *testing.T) {
	w := NewWorld(8)
	acc, err := w.Accountant.RegisterAccount("tmp")
	assert.NoError(t, err)
	nameA, err := w.SpawnRover(acc.Name)
	assert.NoError(t, err)

	infoA, err := w.GetRover(nameA)
	assert.NoError(t, err)

	// Drop a dormant rover on the current position
	infoB := DefaultRover()
	infoB.Name = "abc"
	infoB.Pos = infoA.Pos
	data, err := json.Marshal(infoB)
	assert.NoError(t, err)
	w.Atlas.SetObject(infoA.Pos, Object{Type: roveapi.Object_RoverDormant, Data: data})

	// Enqueue a transfer as well as a dud command
	err = w.Enqueue(nameA,
		&roveapi.Command{Command: roveapi.CommandType_transfer},
		&roveapi.Command{Command: roveapi.CommandType_broadcast, Data: []byte("xyz")})
	assert.NoError(t, err)
	w.Tick()

	// Ensure both command queues are empty
	assert.Empty(t, w.CommandQueue[nameA])
	assert.Empty(t, w.CommandQueue[infoB.Name])

	// Verify the account now controls the new rover
	accountRover, err := w.Accountant.GetValue(acc.Name, "rover")
	assert.NoError(t, err)
	assert.Equal(t, infoB.Name, accountRover)

	// Verify the position now has a dormant rover
	_, obj := w.Atlas.QueryPosition(infoA.Pos)
	assert.Equal(t, roveapi.Object_RoverDormant, obj.Type)

	// Verify the stored data matches
	var stored Rover
	err = json.Unmarshal(obj.Data, &stored)
	assert.NoError(t, err)
	assert.Equal(t, infoA.Name, stored.Name)

	// Verify the new rover data matches what we put in
	infoB2, err := w.GetRover(infoB.Name)
	assert.NoError(t, err)
	assert.Equal(t, infoB.Name, infoB2.Name)

}

func TestCommand_Invalid(t *testing.T) {
	w := NewWorld(8)
	name, err := w.SpawnRover("tmp")
	assert.NoError(t, err)

	err = w.Enqueue(name, &roveapi.Command{Command: roveapi.CommandType_none})
	assert.Error(t, err)
}
