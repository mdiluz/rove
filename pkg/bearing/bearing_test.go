package bearing

import (
	"testing"

	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestDirection(t *testing.T) {
	dir := North

	assert.Equal(t, "North", dir.String())
	assert.Equal(t, "N", dir.ShortString())
	assert.Equal(t, vector.Vector{X: 0, Y: 1}, dir.Vector())

	dir, err := DirectionFromString("N")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = DirectionFromString("n")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = DirectionFromString("north")
	assert.NoError(t, err)
	assert.Equal(t, North, dir)

	dir, err = DirectionFromString("NorthWest")
	assert.NoError(t, err)
	assert.Equal(t, NorthWest, dir)
}
