package game

import (
	"testing"

	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/bearing"
	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	// Very basic for now, nothing to verify
	world := NewWorld(4, 4)
	if world == nil {
		t.Error("Failed to create world")
	}
}

func TestWorld_CreateRover(t *testing.T) {
	world := NewWorld(2, 8)
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
	world := NewWorld(2, 4)
	a, err := world.SpawnRover()
	assert.NoError(t, err)

	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	assert.NotZero(t, attribs.Range, "Rover should not be spawned blind")
	assert.NotZero(t, attribs.Speed, "Rover should not be spawned unable to move")
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld(4, 1)
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
	world := NewWorld(4, 4)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	pos := vector.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newAttribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly set position for rover")

	b := bearing.North
	duration := 1
	newAttribs, err = world.MoveRover(a, b)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(vector.Vector{X: 0, Y: attribs.Speed * duration}) // We should have move one unit of the speed north
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly move position for rover")

	// Place a tile in front of the rover
	assert.NoError(t, world.Atlas.SetTile(vector.Vector{X: 0, Y: 2}, atlas.TileWall))
	newAttribs, err = world.MoveRover(a, b)
	assert.Equal(t, pos, newAttribs.Pos, "Failed to correctly not move position for rover into wall")
}

func TestWorld_RadarFromRover(t *testing.T) {
	// Create world that should have visible walls on the radar
	world := NewWorld(4, 2)
	a, err := world.SpawnRover()
	assert.NoError(t, err)
	b, err := world.SpawnRover()
	assert.NoError(t, err)

	// Set the rover range to a predictable value
	attrib, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	attrib.Range = 4 // Set the range to 4 so we can predict the radar fully
	err = world.SetRoverAttributes(a, attrib)
	assert.NoError(t, err, "Failed to set rover attribs")

	// Warp the rovers into position
	bpos := vector.Vector{X: -3, Y: -3}
	assert.NoError(t, world.WarpRover(b, bpos), "Failed to warp rover")
	assert.NoError(t, world.WarpRover(a, vector.Vector{X: 0, Y: 0}), "Failed to warp rover")

	// Spawn the world wall
	err = world.Atlas.SpawnWalls()
	assert.NoError(t, err)

	radar, err := world.RadarFromRover(a)
	assert.NoError(t, err, "Failed to get radar from rover")
	fullRange := 4 + 4 + 1
	assert.Equal(t, fullRange*fullRange, len(radar), "Radar returned wrong length")

	// It should look like:
	// ---------
	// OOOOOOOO-
	// O------O-
	// O------O-
	// O---R--O-
	// O------O-
	// O------O-
	// OR-----O-
	// OOOOOOOO-
	PrintTiles(radar)

	// Test all expected values
	assert.Equal(t, atlas.TileRover, radar[1+fullRange])
	assert.Equal(t, atlas.TileRover, radar[4+4*fullRange])
	for i := 0; i < 8; i++ {
		assert.Equal(t, atlas.TileWall, radar[i])
		assert.Equal(t, atlas.TileWall, radar[i+(7*9)])
		assert.Equal(t, atlas.TileWall, radar[i*9])
		assert.Equal(t, atlas.TileWall, radar[(i*9)+7])
	}
}
