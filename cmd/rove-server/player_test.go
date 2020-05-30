package main

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	a := NewPlayer()
	b := NewPlayer()
	if a.id == b.id {
		t.Error("Player IDs matched")
	}
}
