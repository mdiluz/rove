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

func TestWorld_CreateInstance(t *testing.T) {
	world := NewWorld()
	a := world.CreateInstance()
	b := world.CreateInstance()

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical instances")
	} else if len(world.Instances) != 2 {
		t.Errorf("Incorrect number of instances created")
	}
}

func TestWorld_DestroyInstance(t *testing.T) {
	world := NewWorld()
	a := world.CreateInstance()
	b := world.CreateInstance()

	err := world.DestroyInstance(a)
	assert.NoError(t, err, "Error returned from instance destroy")

	// Basic duplicate check
	if len(world.Instances) != 1 {
		t.Error("Too many instances left in world")
	} else if _, ok := world.Instances[b]; !ok {
		t.Error("Remaining instance is incorrect")
	}
}
