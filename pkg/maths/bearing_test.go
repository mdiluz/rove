package maths

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirection(t *testing.T) {
	dir := North

	assert.Equal(t, "North", dir.String())
	assert.Equal(t, "N", dir.ShortString())
	assert.Equal(t, Vector{X: 0, Y: 1}, dir.Vector())

	dir, err := BearingFromString("N")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = BearingFromString("n")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = BearingFromString("north")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = BearingFromString("NorthWest")
	assert.NoError(t, err)
	assert.Equal(t, NorthWest, dir)
}
