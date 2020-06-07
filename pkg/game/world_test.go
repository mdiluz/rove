package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	// Very basic for now, nothing to verify
	world := NewWorld()
	if world == nil {
		t.Error("Failed to create world")
	}
}

func TestWorld_CreateRover(t *testing.T) {
	world := NewWorld()
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

func TestWorld_RoverAttributes(t *testing.T) {
	world := NewWorld()
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	assert.NotZero(t, attribs.Range, "Rover should not be spawned blind")
	assert.NotZero(t, attribs.Speed, "Rover should not be spawned unable to move")
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld()
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
	world := NewWorld()
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	pos := Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newAttribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly set position for rover")

	bearing := North
	duration := 1
	newAttribs, err = world.MoveRover(a, bearing)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(Vector{0, attribs.Speed * duration}) // We should have move one unit of the speed north
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly move position for rover")

	// Place a tile in front of the rover
	assert.NoError(t, world.Atlas.SetTile(Vector{0, 2}, TileWall))
	newAttribs, err = world.MoveRover(a, bearing)
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly not move position for rover into wall")
}

func TestWorld_RadarFromRover(t *testing.T) {
	world := NewWorld()
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)
	c, err := world.SpawnRover()
	assert.NoError(t, err)

	// Get a's attributes
	attrib, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	// Warp the rovers so a can see b but not c
	bpos := Vector{attrib.Range - 1, 0}
	cpos := Vector{attrib.Range + 1, 0}
	assert.NoError(t, world.WarpRover(a, Vector{0, 0}), "Failed to warp rover")
	assert.NoError(t, world.WarpRover(b, bpos), "Failed to warp rover")
	assert.NoError(t, world.WarpRover(c, cpos), "Failed to warp rover")

	radar, err := world.RadarFromRover(a)
	assert.NoError(t, err, "Failed to get radar from rover")
	assert.Equal(t, 1, len(radar), "Radar returned wrong number of rovers")

	found := false
	for _, blip := range radar {
		if blip.Position == bpos && blip.Tile == TileRover {
			found = true
		}
	}
	assert.True(t, found, "Rover not found on radar in expected position")

}
