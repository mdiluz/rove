package atlas

import (
	"testing"

	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestAtlas_NewAtlas(t *testing.T) {
	a := NewAtlas(1)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ChunkSize)
	assert.Equal(t, 0, len(a.Chunks)) // Should start empty
}

func TestAtlas_toChunk(t *testing.T) {
	a := NewAtlas(1)
	assert.NotNil(t, a)

	// We start empty so we'll look like this
	chunkID := a.toChunk(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, 0, chunkID)

	// Get a tile to spawn the chunks
	a.GetTile(vector.Vector{})

	// Chunks should look like:
	//  2 | 3
	//  -----
	//  0 | 1
	chunkID = a.toChunk(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, 3, chunkID)
	chunkID = a.toChunk(vector.Vector{X: 0, Y: -1})
	assert.Equal(t, 1, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, 0, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -1, Y: 0})
	assert.Equal(t, 2, chunkID)

	a = NewAtlas(2)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.GetTile(vector.Vector{})
	// Chunks should look like:
	// 2 | 3
	// -----
	// 0 | 1
	chunkID = a.toChunk(vector.Vector{X: 1, Y: 1})
	assert.Equal(t, 3, chunkID)
	chunkID = a.toChunk(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, 1, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -2, Y: -2})
	assert.Equal(t, 0, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -2, Y: 1})
	assert.Equal(t, 2, chunkID)

	a = NewAtlas(2)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.GetTile(vector.Vector{X: 0, Y: 3})
	// Chunks should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	chunkID = a.toChunk(vector.Vector{X: 1, Y: 3})
	assert.Equal(t, 14, chunkID)
	chunkID = a.toChunk(vector.Vector{X: 1, Y: -3})
	assert.Equal(t, 2, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, 5, chunkID)
	chunkID = a.toChunk(vector.Vector{X: -2, Y: 2})
	assert.Equal(t, 13, chunkID)
}

func TestAtlas_GetSetTile(t *testing.T) {
	a := NewAtlas(10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	a.SetTile(vector.Vector{X: 0, Y: 0}, 1)
	tile := a.GetTile(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	// Set another tile to 1 and test it
	a.SetTile(vector.Vector{X: 5, Y: -2}, 2)
	tile = a.GetTile(vector.Vector{X: 5, Y: -2})
	assert.Equal(t, byte(2), tile)
}

func TestAtlas_Grown(t *testing.T) {
	// Start with a small example
	a := NewAtlas(2)
	assert.NotNil(t, a)
	assert.Equal(t, 0, len(a.Chunks))

	// Set a few tiles to values
	a.SetTile(vector.Vector{X: 0, Y: 0}, 1)
	a.SetTile(vector.Vector{X: -1, Y: -1}, 2)
	a.SetTile(vector.Vector{X: 1, Y: -2}, 3)

	// Check tile values
	tile := a.GetTile(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile = a.GetTile(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile = a.GetTile(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)

	tile = a.GetTile(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile = a.GetTile(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile = a.GetTile(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)
}
