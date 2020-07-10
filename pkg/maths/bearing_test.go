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

	dir, err := FromString("N")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = FromString("n")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = FromString("north")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = FromString("NorthWest")
	assert.NoError(t, err)
	assert.Equal(t, NorthWest, dir)
}
