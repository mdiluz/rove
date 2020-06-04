package game

import (
	"testing"

	"github.com/google/uuid"
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
	a := uuid.New()
	b := uuid.New()
	assert.NoError(t, world.SpawnRover(a), "Failed to spawn")
	assert.NoError(t, world.SpawnRover(b), "Failed to spawn")

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical rovers")
	} else if len(world.Rovers) != 2 {
		t.Errorf("Incorrect number of rovers created")
	}
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld()
	a := uuid.New()
	b := uuid.New()
	assert.NoError(t, world.SpawnRover(a), "Failed to spawn")
	assert.NoError(t, world.SpawnRover(b), "Failed to spawn")

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
	a := uuid.New()
	assert.NoError(t, world.SpawnRover(a), "Failed to spawn")

	pos := Vector{
		X: 1.0,
		Y: 2.0,
	}

	err := world.SetPosition(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newpos, err := world.GetPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newpos, "Failed to correctly set position for rover")

	newpos, err = world.MovePosition(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(pos)
	assert.Equal(t, pos, newpos, "Failed to correctly move position for rover")
}
