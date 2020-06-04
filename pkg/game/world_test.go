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
	a := world.SpawnRover()
	b := world.SpawnRover()

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical rovers")
	} else if len(world.Rovers) != 2 {
		t.Errorf("Incorrect number of rovers created")
	}
}

func TestWorld_RoverAttributes(t *testing.T) {
	world := NewWorld()
	a := world.SpawnRover()

	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	assert.NotZero(t, attribs.Sight, "Rover should not be spawned blind")
	assert.NotZero(t, attribs.Speed, "Rover should not be spawned unable to move")
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld()
	a := world.SpawnRover()
	b := world.SpawnRover()

	err := world.DestroyRover(a)
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
	a := world.SpawnRover()
	attribs, err := world.RoverAttributes(a)
	assert.NoError(t, err, "Failed to get rover attribs")

	pos := Vector{
		X: 1.0,
		Y: 2.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newpos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newpos, "Failed to correctly set position for rover")

	bearing := 0.0
	duration := 1.0
	newpos, err = world.MoveRover(a, bearing, duration)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(Vector{0, attribs.Speed * float64(duration)}) // We should have move one unit of the speed north
	assert.Equal(t, pos, newpos, "Failed to correctly move position for rover")
}
