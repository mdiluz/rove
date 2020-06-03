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

func TestWorld_CreateInstance(t *testing.T) {
	world := NewWorld()
	a := uuid.New()
	b := uuid.New()
	assert.NoError(t, world.Spawn(a), "Failed to spawn")
	assert.NoError(t, world.Spawn(b), "Failed to spawn")

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical instances")
	} else if len(world.Instances) != 2 {
		t.Errorf("Incorrect number of instances created")
	}
}

func TestWorld_DestroyInstance(t *testing.T) {
	world := NewWorld()
	a := uuid.New()
	b := uuid.New()
	assert.NoError(t, world.Spawn(a), "Failed to spawn")
	assert.NoError(t, world.Spawn(b), "Failed to spawn")

	err := world.DestroyInstance(a)
	assert.NoError(t, err, "Error returned from instance destroy")

	// Basic duplicate check
	if len(world.Instances) != 1 {
		t.Error("Too many instances left in world")
	} else if _, ok := world.Instances[b]; !ok {
		t.Error("Remaining instance is incorrect")
	}
}

func TestWorld_GetSetMovePosition(t *testing.T) {
	world := NewWorld()
	a := uuid.New()
	assert.NoError(t, world.Spawn(a), "Failed to spawn")

	pos := Vector{
		X: 1.0,
		Y: 2.0,
		Z: 3.0,
	}

	err := world.SetPosition(a, pos)
	assert.NoError(t, err, "Failed to set position for instance")

	newpos, err := world.GetPosition(a)
	assert.NoError(t, err, "Failed to set position for instance")
	assert.Equal(t, pos, newpos, "Failed to correctly set position for instance")

	newpos, err = world.MovePosition(a, pos)
	assert.NoError(t, err, "Failed to set position for instance")
	pos.Add(pos)
	assert.Equal(t, pos, newpos, "Failed to correctly move position for instance")
}
