package game

import (
	"fmt"
	"log"
	"math/rand"
)

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []Tile `json:"tiles"`
}

// Atlas represents a grid of Chunks
type Atlas struct {
	// Chunks represents all chunks in the world
	// This is intentionally not a 2D array so it can be expanded in all directions
	Chunks []Chunk `json:"chunks"`

	// size is the current width/height of the given atlas
	Size int `json:"size"`

	// ChunkSize is the dimensions of each chunk
	ChunkSize int `json:"chunksize"`
}

// NewAtlas creates a new empty atlas
func NewAtlas(size int, chunkSize int) Atlas {
	if size%2 != 0 {
		log.Fatal("atlas size must always be even")
	}

	a := Atlas{
		Size:      size,
		Chunks:    make([]Chunk, size*size),
		ChunkSize: chunkSize,
	}

	// Initialise all the chunks
	for i := range a.Chunks {
		a.Chunks[i] = Chunk{
			Tiles: make([]Tile, chunkSize*chunkSize),
		}
	}

	return a
}

// SpawnWorld spawns the current world
func (a *Atlas) SpawnWorld() error {
	extent := a.ChunkSize * (a.Size / 2)

	// Pepper the current world with rocks
	for i := -extent; i < extent; i++ {
		for j := -extent; j < extent; j++ {
			if rand.Intn(16) == 0 {
				if err := a.SetTile(Vector{i, j}, TileRock); err != nil {
					return err
				}
			}
		}
	}

	// Surround the atlas in walls
	for i := -extent; i < extent; i++ {

		if err := a.SetTile(Vector{i, extent - 1}, TileWall); err != nil {
			return err
		} else if a.SetTile(Vector{extent - 1, i}, TileWall); err != nil {
			return err
		} else if a.SetTile(Vector{-extent, i}, TileWall); err != nil {
			return err
		} else if a.SetTile(Vector{i, extent - 1}, TileWall); err != nil {
			return err
		}
	}

	return nil
}

// SetTile sets an individual tile's kind
func (a *Atlas) SetTile(v Vector, tile Tile) error {
	chunk := a.ToChunk(v)
	if chunk >= len(a.Chunks) {
		return fmt.Errorf("location outside of allocated atlas")
	}

	local := a.ToChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize
	if tileId >= len(a.Chunks[chunk].Tiles) {
		return fmt.Errorf("location outside of allocated chunk")
	}
	a.Chunks[chunk].Tiles[tileId] = tile
	return nil
}

// GetTile will return an individual tile
func (a *Atlas) GetTile(v Vector) (Tile, error) {
	chunk := a.ToChunk(v)
	if chunk >= len(a.Chunks) {
		return 0, fmt.Errorf("location outside of allocated atlas")
	}

	local := a.ToChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize
	if tileId >= len(a.Chunks[chunk].Tiles) {
		return 0, fmt.Errorf("location outside of allocated chunk")
	}

	return a.Chunks[chunk].Tiles[tileId], nil
}

// ToChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) ToChunkLocal(v Vector) Vector {
	return Vector{Pmod(v.X, a.ChunkSize), Pmod(v.Y, a.ChunkSize)}
}

// GetChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) ToWorld(local Vector, chunk int) Vector {
	return a.ChunkOrigin(chunk).Added(local)
}

// GetChunkID gets the chunk ID for a position in the world
func (a *Atlas) ToChunk(v Vector) int {
	local := a.ToChunkLocal(v)
	// Get the chunk origin itself
	origin := v.Added(local.Negated())
	// Divided it by the number of chunks
	origin = origin.Divided(a.ChunkSize)
	// Shift it by our size (our origin is in the middle)
	origin = origin.Added(Vector{a.Size / 2, a.Size / 2})
	// Get the ID based on the final values
	return (a.Size * origin.Y) + origin.X
}

// ChunkOrigin gets the chunk origin for a given chunk index
func (a *Atlas) ChunkOrigin(chunk int) Vector {
	v := Vector{
		X: Pmod(chunk, a.Size) - (a.Size / 2),
		Y: (chunk / a.Size) - (a.Size / 2),
	}

	return v.Multiplied(a.ChunkSize)
}

// GetWorldExtent gets the min and max valid coordinates of world
func (a *Atlas) GetWorldExtents() (min Vector, max Vector) {
	min = Vector{
		-(a.Size / 2) * a.ChunkSize,
		-(a.Size / 2) * a.ChunkSize,
	}
	max = Vector{
		-min.X - 1,
		-min.Y - 1,
	}
	return
}

// Grow will return a grown copy of the current atlas
func (a *Atlas) Grow(size int) error {
	if size%2 != 0 {
		return fmt.Errorf("atlas size must always be even")
	}
	delta := size - a.Size
	if delta < 0 {
		return fmt.Errorf("Cannot shrink an atlas")
	} else if delta == 0 {
		return nil
	}

	// Create a new atlas
	newAtlas := NewAtlas(size, a.ChunkSize)

	// Copy old chunks into new chunks
	for index, chunk := range a.Chunks {
		// Calculate the new chunk location and copy over the data
		newAtlas.Chunks[newAtlas.ToChunk(a.ChunkOrigin(index))] = chunk
	}

	// Copy the new atlas data into this one
	*a = newAtlas

	// Return the new atlas
	return nil
}
