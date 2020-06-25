package atlas

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/vector"
)

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []byte `json:"tiles"`
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
func NewAtlas(size, chunkSize int) Atlas {
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
		tiles := make([]byte, chunkSize*chunkSize)
		for i := 0; i < len(tiles); i++ {
			tiles[i] = TileEmpty
		}
		a.Chunks[i] = Chunk{
			Tiles: tiles,
		}
	}

	return a
}

// SpawnRocks peppers the world with rocks
func (a *Atlas) SpawnRocks() error {
	extent := a.ChunkSize * (a.Size / 2)

	// Pepper the current world with rocks
	for i := -extent; i < extent; i++ {
		for j := -extent; j < extent; j++ {
			if rand.Intn(16) == 0 {
				if err := a.SetTile(vector.Vector{X: i, Y: j}, TileSmallRock); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SpawnWalls spawns the around the world
func (a *Atlas) SpawnWalls() error {
	extent := a.ChunkSize * (a.Size / 2)

	// Surround the atlas in walls
	for i := -extent; i < extent; i++ {

		if err := a.SetTile(vector.Vector{X: i, Y: extent - 1}, TileLargeRock); err != nil { // N
			return err
		} else if err := a.SetTile(vector.Vector{X: extent - 1, Y: i}, TileLargeRock); err != nil { // E
			return err
		} else if err := a.SetTile(vector.Vector{X: i, Y: -extent}, TileLargeRock); err != nil { // S
			return err
		} else if err := a.SetTile(vector.Vector{X: -extent, Y: i}, TileLargeRock); err != nil { // W
			return err
		}
	}

	return nil
}

// SetTile sets an individual tile's kind
func (a *Atlas) SetTile(v vector.Vector, tile byte) error {
	chunk := a.toChunk(v)
	if chunk >= len(a.Chunks) {
		return fmt.Errorf("location outside of allocated atlas")
	}

	local := a.toChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize
	if tileId >= len(a.Chunks[chunk].Tiles) {
		return fmt.Errorf("location outside of allocated chunk")
	}
	a.Chunks[chunk].Tiles[tileId] = tile
	return nil
}

// GetTile will return an individual tile
func (a *Atlas) GetTile(v vector.Vector) (byte, error) {
	chunk := a.toChunk(v)
	if chunk >= len(a.Chunks) {
		return 0, fmt.Errorf("location outside of allocated atlas")
	}

	local := a.toChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize
	if tileId >= len(a.Chunks[chunk].Tiles) {
		return 0, fmt.Errorf("location outside of allocated chunk")
	}

	return a.Chunks[chunk].Tiles[tileId], nil
}

// toChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) toChunkLocal(v vector.Vector) vector.Vector {
	return vector.Vector{X: maths.Pmod(v.X, a.ChunkSize), Y: maths.Pmod(v.Y, a.ChunkSize)}
}

// GetChunkID gets the chunk ID for a position in the world
func (a *Atlas) toChunk(v vector.Vector) int {
	local := a.toChunkLocal(v)
	// Get the chunk origin itself
	origin := v.Added(local.Negated())
	// Divided it by the number of chunks
	origin = origin.Divided(a.ChunkSize)
	// Shift it by our size (our origin is in the middle)
	origin = origin.Added(vector.Vector{X: a.Size / 2, Y: a.Size / 2})
	// Get the ID based on the final values
	return (a.Size * origin.Y) + origin.X
}

// chunkOrigin gets the chunk origin for a given chunk index
func (a *Atlas) chunkOrigin(chunk int) vector.Vector {
	v := vector.Vector{
		X: maths.Pmod(chunk, a.Size) - (a.Size / 2),
		Y: (chunk / a.Size) - (a.Size / 2),
	}

	return v.Multiplied(a.ChunkSize)
}

// GetWorldExtent gets the min and max valid coordinates of world
func (a *Atlas) GetWorldExtents() (min, max vector.Vector) {
	min = vector.Vector{
		X: -(a.Size / 2) * a.ChunkSize,
		Y: -(a.Size / 2) * a.ChunkSize,
	}
	max = vector.Vector{
		X: -min.X - 1,
		Y: -min.Y - 1,
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
		return fmt.Errorf("cannot shrink an atlas")
	} else if delta == 0 {
		return nil
	}

	// Create a new atlas
	newAtlas := NewAtlas(size, a.ChunkSize)

	// Copy old chunks into new chunks
	for index, chunk := range a.Chunks {
		// Calculate the new chunk location and copy over the data
		newAtlas.Chunks[newAtlas.toChunk(a.chunkOrigin(index))] = chunk
	}

	// Copy the new atlas data into this one
	*a = newAtlas

	// Return the new atlas
	return nil
}
