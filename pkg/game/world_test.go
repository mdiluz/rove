package game

import (
	"testing"

	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/bearing"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
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

	pos := vector.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newPos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newPos, "Failed to correctly set position for rover")

	b := bearing.North
	newPos, err = world.MoveRover(a, b)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(vector.Vector{X: 0, Y: 1})
	assert.Equal(t, pos, newPos, "Failed to correctly move position for rover")

	// Place a tile in front of the rover
	world.Atlas.SetObject(vector.Vector{X: 0, Y: 2}, objects.Object{Type: objects.LargeRock})
	newPos, err = world.MoveRover(a, b)
	assert.Equal(t, pos, newPos, "Failed to correctly not move position for rover into wall")
}

func TestWorld_RadarFromRover(t *testing.T) {
	// Create world that should have visible walls on the radar
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)

	// Warp the rovers into position
	bpos := vector.Vector{X: -3, Y: -3}
	assert.NoError(t, world.WarpRover(b, bpos), "Failed to warp rover")
	assert.NoError(t, world.WarpRover(a, vector.Vector{X: 0, Y: 0}), "Failed to warp rover")

	radar, objs, err := world.RadarFromRover(a)
	assert.NoError(t, err, "Failed to get radar from rover")
	fullRange := 4 + 4 + 1
	assert.Equal(t, fullRange*fullRange, len(radar), "Radar returned wrong length")
	assert.Equal(t, fullRange*fullRange, len(objs), "Radar returned wrong length")

	// Test the expected values
	assert.Equal(t, byte(objects.Rover), objs[1+fullRange])
	assert.Equal(t, byte(objects.Rover), objs[4+4*fullRange])
}

func TestWorld_RoverStash(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := vector.Vector{
		X: 0.0,
		Y: 0.0,
	}

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
}

func TestWorld_RoverDamage(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := vector.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	info, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")

	world.Atlas.SetObject(vector.Vector{X: 0.0, Y: 1.0}, objects.Object{Type: objects.LargeRock})

	vec, err := world.MoveRover(a, bearing.North)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, vec, "Rover managed to move into large rock")

	newinfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, info.Integrity-1, newinfo.Integrity, "rover should have lost integrity")
}

func TestWorld_RoverRepair(t *testing.T) {
	world := NewWorld(2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	pos := vector.Vector{
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

	world.Atlas.SetObject(vector.Vector{X: 0.0, Y: 1.0}, objects.Object{Type: objects.LargeRock})

	// Try and bump into the rock
	vec, err := world.MoveRover(a, bearing.North)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, vec, "Rover managed to move into large rock")

	newinfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity-1, newinfo.Integrity, "rover should have lost integrity")

	err = world.ExecuteCommand(&Command{Command: CommandRepair}, a)
	assert.NoError(t, err, "Failed to repair rover")

	newinfo, err = world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity, newinfo.Integrity, "rover should have gained integrity")

	// Check again that it can't repair past the max
	world.Atlas.SetObject(pos, objects.Object{Type: objects.SmallRock})
	o, err = world.RoverStash(a)
	assert.NoError(t, err, "Failed to stash")
	assert.Equal(t, objects.SmallRock, o, "Failed to get correct object")

	err = world.ExecuteCommand(&Command{Command: CommandRepair}, a)
	assert.NoError(t, err, "Failed to repair rover")

	newinfo, err = world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, originalInfo.Integrity, newinfo.Integrity, "rover should have kept the same integrity")
}
