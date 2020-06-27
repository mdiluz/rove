package atlas

import (
	"log"
	"math/rand"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
)

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []byte `json:"tiles"`
}

// SpawnContent will create a chunk and fill it with spawned tiles
func (c *Chunk) SpawnContent(size int) {
	c.Tiles = make([]byte, size*size)
	for i := 0; i < len(c.Tiles); i++ {
		c.Tiles[i] = objects.Empty
	}

	// For now, fill it randomly with objects
	for i := range c.Tiles {
		if rand.Intn(16) == 0 {
			c.Tiles[i] = objects.LargeRock
		} else if rand.Intn(32) == 0 {
			c.Tiles[i] = objects.SmallRock
		}
	}
}

// Atlas represents a grid of Chunks
type Atlas struct {
	// Chunks represents all chunks in the world
	// This is intentionally not a 2D array so it can be expanded in all directions
	Chunks []Chunk `json:"chunks"`

	// CurrentSize is the current width/height of the given atlas
	CurrentSize int `json:"currentSize"`

	// ChunkSize is the dimensions of each chunk
	ChunkSize int `json:"chunksize"`
}

// NewAtlas creates a new empty atlas
func NewAtlas(chunkSize int) Atlas {
	return Atlas{
		CurrentSize: 0,
		Chunks:      nil,
		ChunkSize:   chunkSize,
	}
}

// SetTile sets an individual tile's kind
func (a *Atlas) SetTile(v vector.Vector, tile byte) {
	// Get the chunk, expand, and spawn it if needed
	c := a.toChunkWithGrow(v)
	chunk := a.Chunks[c]
	if chunk.Tiles == nil {
		chunk.SpawnContent(a.ChunkSize)
	}

	local := a.toChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize

	// Sanity check
	if tileId >= len(chunk.Tiles) || tileId < 0 {
		log.Fatalf("Local tileID is not in valid chunk, somehow, this means something is very wrong")
	}

	// Set the chunk back
	chunk.Tiles[tileId] = tile
	a.Chunks[c] = chunk
}

// GetTile will return an individual tile
func (a *Atlas) GetTile(v vector.Vector) byte {
	// Get the chunk, expand, and spawn it if needed
	c := a.toChunkWithGrow(v)
	chunk := a.Chunks[c]
	if chunk.Tiles == nil {
		chunk.SpawnContent(a.ChunkSize)
	}

	local := a.toChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize

	// Sanity check
	if tileId >= len(chunk.Tiles) || tileId < 0 {
		log.Fatalf("Local tileID is not in valid chunk, somehow, this means something is very wrong")
	}

	return chunk.Tiles[tileId]
}

// toChunkWithGrow will expand the atlas for a given tile, returns the new chunk
func (a *Atlas) toChunkWithGrow(v vector.Vector) int {
	for {
		// Get the chunk, and grow looping until we have a valid chunk
		chunk := a.toChunk(v)
		if chunk >= len(a.Chunks) || chunk < 0 {
			a.grow()
		} else {
			return chunk
		}
	}
}

// toChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) toChunkLocal(v vector.Vector) vector.Vector {
	return vector.Vector{X: maths.Pmod(v.X, a.ChunkSize), Y: maths.Pmod(v.Y, a.ChunkSize)}
}

// GetChunkID gets the current chunk ID for a position in the world
func (a *Atlas) toChunk(v vector.Vector) int {
	local := a.toChunkLocal(v)
	// Get the chunk origin itself
	origin := v.Added(local.Negated())
	// Divided it by the number of chunks
	origin = origin.Divided(a.ChunkSize)
	// Shift it by our size (our origin is in the middle)
	origin = origin.Added(vector.Vector{X: a.CurrentSize / 2, Y: a.CurrentSize / 2})
	// Get the ID based on the final values
	return (a.CurrentSize * origin.Y) + origin.X
}

// chunkOrigin gets the chunk origin for a given chunk index
func (a *Atlas) chunkOrigin(chunk int) vector.Vector {
	v := vector.Vector{
		X: maths.Pmod(chunk, a.CurrentSize) - (a.CurrentSize / 2),
		Y: (chunk / a.CurrentSize) - (a.CurrentSize / 2),
	}

	return v.Multiplied(a.ChunkSize)
}

// grow will expand the current atlas in all directions by one chunk
func (a *Atlas) grow() error {
	// Create a new atlas
	newAtlas := NewAtlas(a.ChunkSize)

	// Expand by one on each axis
	newAtlas.CurrentSize = a.CurrentSize + 2

	// Allocate the new atlas chunks
	// These chunks will have nil tile slices
	newAtlas.Chunks = make([]Chunk, newAtlas.CurrentSize*newAtlas.CurrentSize)

	// Copy all old chunks into the new atlas
	for index, chunk := range a.Chunks {
		// Calculate the new chunk location and copy over the data
		newAtlas.Chunks[newAtlas.toChunk(a.chunkOrigin(index))] = chunk
	}

	// Copy the new atlas data into this one
	*a = newAtlas

	// Return the new atlas
	return nil
}
