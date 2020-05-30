package rovegame

import (
	"testing"
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
	}
}
