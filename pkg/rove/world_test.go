package rove

import (
	"testing"

	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/roveapi"
	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	// Very basic for now, nothing to verify
	world := NewWorld(4)
	if world == nil {
		t.Error("Failed to create world")
	}
}

func TestWorld_CreateRover(t *testing.T) {
	world := NewWorld(8)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical rovers")
	} else if len(world.Rovers) != 2 {
		t.Errorf("Incorrect number of rovers created")
	}
}

func TestWorld_GetRover(t *testing.T) {
	world := NewWorld(4)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	assert.NotZero(t, rover.Range, "Rover should not be spawned blind")
	assert.Contains(t, rover.Logs[len(rover.Logs)-1].Text, "created", "Rover logs should contain the creation")
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld(1)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)

	err = world.DestroyRover(a)
	assert.NoError(t, err, "Error returned from rover destroy")

	// Basic duplicate check
	if len(world.Rovers) != 1 {
		t.Error("Too many rovers left in world")
	} else if _, ok := world.Rovers[b]; !ok {
		t.Error("Remaining rover is incorrect")
	}
}

func TestWorld_GetSetMovePosition(t *testing.T) {
	world := NewWorld(4)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newPos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newPos, "Failed to correctly set position for rover")

	b := maths.North
	newPos, err = world.MoveRover(a, b)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(maths.Vector{X: 0, Y: 1})
	assert.Equal(t, pos, newPos, "Failed to correctly move position for rover")

	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover information")
	assert.Equal(t, rover.MaximumCharge-1, rover.Charge, "Rover should have lost charge for moving")
	assert.Contains(t, rover.Logs[len(rover.Logs)-1].Text, "moved", "Rover logs should contain the move")

	// Place a tile in front of the rover
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 2}, objects.Object{Type: objects.LargeRock})
	newPos, err = world.MoveRover(a, b)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, newPos, "Failed to correctly not move position for rover into wall")

	rover, err = world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover information")
	assert.Equal(t, rover.MaximumCharge-2, rover.Charge, "Rover should have lost charge for move attempt")
}

func TestWorld_RadarFromRover(t *testing.T) {
	// Create world that should have visible walls on the radar
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)

	// Warp the rovers into position
	bpos := maths.Vector{X: -3, Y: -3}
	assert.NoError(t, world.WarpRover(b, bpos), "Failed to warp rover")
	assert.NoError(t, world.WarpRover(a, maths.Vector{X: 0, Y: 0}), "Failed to warp rover")

	radar, objs, err := world.RadarFromRover(a)
	assert.NoError(t, err, "Failed to get radar from rover")
	fullRange := 4 + 4 + 1
	assert.Equal(t, fullRange*fullRange, len(radar), "Radar returned wrong length")
	assert.Equal(t, fullRange*fullRange, len(objs), "Radar returned wrong length")

	// Test the expected values
	assert.Equal(t, byte(objects.Rover), objs[1+fullRange])
	assert.Equal(t, byte(objects.Rover), objs[4+4*fullRange])

	// Check the radar results are stable
	radar1, objs1, err := world.RadarFromRover(a)
	assert.NoError(t, err)
	radar2, objs2, err := world.RadarFromRover(a)
	assert.NoError(t, err)
	assert.Equal(t, radar1, radar2)
	assert.Equal(t, objs1, objs2)
}

func TestWorld_RoverStash(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	world.Atlas.SetObject(pos, objects.Object{Type: objects.None})
	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	// Set to a traversible tile
	world.Atlas.SetTile(pos, atlas.TileRock)

	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover")

	for i := 0; i < rover.Capacity; i++ {
		// Place an object
		world.Atlas.SetObject(pos, objects.Object{Type: objects.SmallRock})

		// Pick it up
		o, err := world.RoverStash(a)
		assert.NoError(t, err, "Failed to stash")
		assert.Equal(t, objects.SmallRock, o, "Failed to get correct object")

		// Check it's gone
		_, obj := world.Atlas.QueryPosition(pos)
		assert.Equal(t, objects.None, obj.Type, "Stash failed to remove object from atlas")

		// Check we have it
		inv, err := world.RoverInventory(a)
		assert.NoError(t, err, "Failed to get inventory")
		assert.Equal(t, i+1, len(inv))
		assert.Equal(t, objects.Object{Type: objects.SmallRock}, inv[i])

		// Check that this did reduce the charge
		info, err := world.GetRover(a)
		assert.NoError(t, err, "Failed to get rover")
		assert.Equal(t, info.MaximumCharge-(i+1), info.Charge, "Rover lost charge for stash")
		assert.Contains(t, info.Logs[len(info.Logs)-1].Text, "stashed", "Rover logs should contain the move")
	}

	// Recharge the rover
	for i := 0; i < rover.MaximumCharge; i++ {
		_, err = world.RoverRecharge(a)
		assert.NoError(t, err)

	}

	// Place an object
	world.Atlas.SetObject(pos, objects.Object{Type: objects.SmallRock})

	// Try to pick it up
	o, err := world.RoverStash(a)
	assert.NoError(t, err, "Failed to stash")
	assert.Equal(t, objects.None, o, "Failed to get correct object")

	// Check it's still there
	_, obj := world.Atlas.QueryPosition(pos)
	assert.Equal(t, objects.SmallRock, obj.Type, "Stash failed to remove object from atlas")

	// Check we don't have it
	inv, err := world.RoverInventory(a)
	assert.NoError(t, err, "Failed to get inventory")
	assert.Equal(t, rover.Capacity, len(inv))

	// Check that this didn't reduce the charge
	info, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover")
	assert.Equal(t, info.MaximumCharge, info.Charge, "Rover lost charge for non-stash")
}

func TestWorld_RoverDamage(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	info, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")

	world.Atlas.SetObject(maths.Vector{X: 0.0, Y: 1.0}, objects.Object{Type: objects.LargeRock})

	vec, err := world.MoveRover(a, maths.North)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, vec, "Rover managed to move into large rock")

	newinfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, info.Integrity-1, newinfo.Integrity, "rover should have lost integrity")
	assert.Contains(t, newinfo.Logs[len(newinfo.Logs)-1].Text, "collision", "Rover logs should contain the collision")
}

func TestWorld_RoverRepair(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	world.Atlas.SetTile(pos, atlas.TileNone)
	world.Atlas.SetObject(pos, objects.Object{Type: objects.None})

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	originalInfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")

	// Pick up something to repair with
	world.Atlas.SetObject(pos, objects.Object{Type: objects.SmallRock})
	o, err := world.RoverStash(a)
	assert.NoError(t, err, "Failed to stash")
	assert.Equal(t, objects.SmallRock, o, "Failed to get correct object")

	world.Atlas.SetObject(maths.Vector{X: 0.0, Y: 1.0}, objects.Object{Type: objects.LargeRock})

	// Try and bump into the rock
	vec, err := world.MoveRover(a, maths.North)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, vec, "Rover managed to move into large rock")

	newinfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity-1, newinfo.Integrity, "rover should have lost integrity")

	err = world.ExecuteCommand(&Command{Command: roveapi.CommandType_repair}, a)
	assert.NoError(t, err, "Failed to repair rover")

	newinfo, err = world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity, newinfo.Integrity, "rover should have gained integrity")
	assert.Contains(t, newinfo.Logs[len(newinfo.Logs)-1].Text, "repair", "Rover logs should contain the repair")

	// Check again that it can't repair past the max
	world.Atlas.SetObject(pos, objects.Object{Type: objects.SmallRock})
	o, err = world.RoverStash(a)
	assert.NoError(t, err, "Failed to stash")
	assert.Equal(t, objects.SmallRock, o, "Failed to get correct object")

	err = world.ExecuteCommand(&Command{Command: roveapi.CommandType_repair}, a)
	assert.NoError(t, err, "Failed to repair rover")

	newinfo, err = world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity, newinfo.Integrity, "rover should have kept the same integrity")
}

func TestWorld_Charge(t *testing.T) {
	world := NewWorld(4)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	// Get the rover information
	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover information")
	assert.Equal(t, rover.MaximumCharge, rover.Charge, "Rover should start with maximum charge")

	// Use up all the charge
	for i := 0; i < rover.MaximumCharge; i++ {
		// Get the initial position
		initialPos, err := world.RoverPosition(a)
		assert.NoError(t, err, "Failed to get position for rover")

		// Ensure the path ahead is empty
		world.Atlas.SetTile(initialPos.Added(maths.North.Vector()), atlas.TileRock)
		world.Atlas.SetObject(initialPos.Added(maths.North.Vector()), objects.Object{Type: objects.None})

		// Try and move north (along unblocked path)
		newPos, err := world.MoveRover(a, maths.North)
		assert.NoError(t, err, "Failed to set position for rover")
		assert.Equal(t, initialPos.Added(maths.North.Vector()), newPos, "Failed to correctly move position for rover")

		// Ensure rover lost charge
		rover, err := world.GetRover(a)
		assert.NoError(t, err, "Failed to get rover information")
		assert.Equal(t, rover.MaximumCharge-(i+1), rover.Charge, "Rover should have lost charge")
	}

}

func TestWorld_Daytime(t *testing.T) {
	world := NewWorld(1)

	a, err := world.SpawnRover()
	assert.NoError(t, err)

	// Remove rover charge
	rover := world.Rovers[a]
	rover.Charge = 0
	world.Rovers[a] = rover

	// Try and recharge, should work
	_, err = world.RoverRecharge(a)
	assert.NoError(t, err)
	assert.Equal(t, 1, world.Rovers[a].Charge)

	// Loop for half the day
	for i := 0; i < world.TicksPerDay/2; i++ {
		assert.True(t, world.Daytime())
		world.ExecuteCommandQueues()
	}

	// Remove rover charge again
	rover = world.Rovers[a]
	rover.Charge = 0
	world.Rovers[a] = rover

	// Try and recharge, should fail
	_, err = world.RoverRecharge(a)
	assert.NoError(t, err)
	assert.Equal(t, 0, world.Rovers[a].Charge)

	// Loop for half the day
	for i := 0; i < world.TicksPerDay/2; i++ {
		assert.False(t, world.Daytime())
		world.ExecuteCommandQueues()
	}
}

func TestWorld_Broadcast(t *testing.T) {
	world := NewWorld(8)

	a, err := world.SpawnRover()
	assert.NoError(t, err)

	b, err := world.SpawnRover()
	assert.NoError(t, err)

	// Warp rovers near to eachother
	assert.NoError(t, world.WarpRover(a, maths.Vector{X: 0, Y: 0}))
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: 1, Y: 0}))

	// Broadcast from a
	assert.NoError(t, world.RoverBroadcast(a, []byte{'A', 'B', 'C'}))

	// Check if b heard it
	ra, err := world.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, ra.MaximumCharge-1, ra.Charge, "A should have used a charge to broadcast")
	assert.Contains(t, ra.Logs[len(ra.Logs)-1].Text, "ABC", "Rover B should have heard the broadcast")

	// Check if a logged it
	rb, err := world.GetRover(b)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "ABC", "Rover A should have logged it's broadcast")

	// Warp B outside of the range of A
	world.Atlas.SetObject(maths.Vector{X: ra.Range, Y: 0}, objects.Object{Type: objects.None})
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: ra.Range, Y: 0}))

	// Broadcast from a again
	assert.NoError(t, world.RoverBroadcast(a, []byte{'X', 'Y', 'Z'}))

	// Check if b heard it
	ra, err = world.GetRover(b)
	assert.NoError(t, err)
	assert.NotContains(t, ra.Logs[len(ra.Logs)-1].Text, "XYZ", "Rover B should not have heard the broadcast")

	// Check if a logged it
	rb, err = world.GetRover(a)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "XYZ", "Rover A should have logged it's broadcast")

	// Warp B outside of the range of A
	world.Atlas.SetObject(maths.Vector{X: ra.Range + 1, Y: 0}, objects.Object{Type: objects.None})
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: ra.Range + 1, Y: 0}))

	// Broadcast from a again
	assert.NoError(t, world.RoverBroadcast(a, []byte{'H', 'J', 'K'}))

	// Check if b heard it
	ra, err = world.GetRover(b)
	assert.NoError(t, err)
	assert.NotContains(t, ra.Logs[len(ra.Logs)-1].Text, "HJK", "Rover B should have heard the broadcast")

	// Check if a logged it
	rb, err = world.GetRover(a)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "HJK", "Rover A should have logged it's broadcast")
}
