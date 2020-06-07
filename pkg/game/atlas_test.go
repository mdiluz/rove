package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtlas_NewAtlas(t *testing.T) {
	a := NewAtlas(2, 1)
	assert.NotNil(t, a)
	// Tiles should look like: 2 | 3
	//  -----
	//  0 | 1
	assert.Equal(t, 4, len(a.Chunks))

	a = NewAtlas(4, 1)
	assert.NotNil(t, a)
	// Tiles should look like: 2 | 3
	//  -----
	//  0 | 1
	assert.Equal(t, 16, len(a.Chunks))
}

func TestAtlas_ToChunk(t *testing.T) {
	a := NewAtlas(2, 1)
	assert.NotNil(t, a)
	// Tiles should look like: 2 | 3
	//  -----
	//  0 | 1
	tile := a.ToChunk(Vector{0, 0})
	assert.Equal(t, 3, tile)
	tile = a.ToChunk(Vector{0, -1})
	assert.Equal(t, 1, tile)
	tile = a.ToChunk(Vector{-1, -1})
	assert.Equal(t, 0, tile)
	tile = a.ToChunk(Vector{-1, 0})
	assert.Equal(t, 2, tile)

	a = NewAtlas(2, 2)
	assert.NotNil(t, a)
	// Tiles should look like:
	// 2 | 3
	// -----
	// 0 | 1
	tile = a.ToChunk(Vector{1, 1})
	assert.Equal(t, 3, tile)
	tile = a.ToChunk(Vector{1, -2})
	assert.Equal(t, 1, tile)
	tile = a.ToChunk(Vector{-2, -2})
	assert.Equal(t, 0, tile)
	tile = a.ToChunk(Vector{-2, 1})
	assert.Equal(t, 2, tile)

	a = NewAtlas(4, 2)
	assert.NotNil(t, a)
	// Tiles should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	tile = a.ToChunk(Vector{1, 3})
	assert.Equal(t, 14, tile)
	tile = a.ToChunk(Vector{1, -3})
	assert.Equal(t, 2, tile)
	tile = a.ToChunk(Vector{-1, -1})
	assert.Equal(t, 5, tile)
	tile = a.ToChunk(Vector{-2, 2})
	assert.Equal(t, 13, tile)
}

func TestAtlas_GetSetTile(t *testing.T) {
	a := NewAtlas(4, 10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	assert.NoError(t, a.SetTile(Vector{0, 0}, 1))
	tile, err := a.GetTile(Vector{0, 0})
	assert.NoError(t, err)
	assert.Equal(t, Tile(1), tile)

	// Set another tile to 1 and test it
	assert.NoError(t, a.SetTile(Vector{5, -2}, 2))
	tile, err = a.GetTile(Vector{5, -2})
	assert.NoError(t, err)
	assert.Equal(t, Tile(2), tile)
}

func TestAtlas_Grown(t *testing.T) {
	// Start with a small example
	a := NewAtlas(2, 2)
	assert.NotNil(t, a)
	assert.Equal(t, 4, len(a.Chunks))

	// Set a few tiles to values
	assert.NoError(t, a.SetTile(Vector{0, 0}, 1))
	assert.NoError(t, a.SetTile(Vector{-1, -1}, 2))
	assert.NoError(t, a.SetTile(Vector{1, -2}, 3))

	// Grow once to just double it
	err := a.Grow(4)
	assert.NoError(t, err)
	assert.Equal(t, 16, len(a.Chunks))

	tile, err := a.GetTile(Vector{0, 0})
	assert.NoError(t, err)
	assert.Equal(t, Tile(1), tile)

	tile, err = a.GetTile(Vector{-1, -1})
	assert.NoError(t, err)
	assert.Equal(t, Tile(2), tile)

	tile, err = a.GetTile(Vector{1, -2})
	assert.NoError(t, err)
	assert.Equal(t, Tile(3), tile)

	// Grow it again even bigger
	err = a.Grow(10)
	assert.NoError(t, err)
	assert.Equal(t, 100, len(a.Chunks))

	tile, err = a.GetTile(Vector{0, 0})
	assert.NoError(t, err)
	assert.Equal(t, Tile(1), tile)

	tile, err = a.GetTile(Vector{-1, -1})
	assert.NoError(t, err)
	assert.Equal(t, Tile(2), tile)

	tile, err = a.GetTile(Vector{1, -2})
	assert.NoError(t, err)
	assert.Equal(t, Tile(3), tile)
}

func TestAtlas_SpawnWorld(t *testing.T) {
	// Start with a small example
	a := NewAtlas(2, 2)
	assert.NotNil(t, a)
	assert.Equal(t, 4, len(a.Chunks))

	assert.NoError(t, a.SpawnWorld())
	tile, err := a.GetTile(Vector{0, 0})
	assert.NoError(t, err)
	assert.Equal(t, TileEmpty, tile)

	tile, err = a.GetTile(Vector{1, 1})
	assert.NoError(t, err)
	assert.Equal(t, TileWall, tile)

	tile, err = a.GetTile(Vector{-1, -1})
	assert.NoError(t, err)
	assert.Equal(t, TileEmpty, tile)

	tile, err = a.GetTile(Vector{-2, -2})
	assert.NoError(t, err)
	assert.Equal(t, TileWall, tile)
}
