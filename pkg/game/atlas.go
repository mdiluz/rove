package game

import "fmt"

// Kind represents the type of a tile on the map
type Kind byte

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []Kind `json:"tiles"`
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
func NewAtlas(radius int, chunkSize int) Atlas {
	size := radius * 2
	a := Atlas{
		Size:      size,
		Chunks:    make([]Chunk, size*size),
		ChunkSize: chunkSize,
	}

	// Initialise all the chunks
	for i := range a.Chunks {
		a.Chunks[i] = Chunk{
			Tiles: make([]Kind, chunkSize*chunkSize),
		}
	}

	return a
}

// SetTile sets an individual tile's kind
func (a *Atlas) SetTile(v Vector, tile Kind) error {
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
func (a *Atlas) GetTile(v Vector) (Kind, error) {
	chunk := a.ToChunk(v)
	if chunk > len(a.Chunks) {
		return 0, fmt.Errorf("location outside of allocated atlas")
	}

	local := a.ToChunkLocal(v)
	tileId := local.X + local.Y*a.ChunkSize
	if tileId > len(a.Chunks[chunk].Tiles) {
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

// Grown will return a grown copy of the current atlas
func (a *Atlas) Grown(newRadius int) (Atlas, error) {
	delta := (newRadius * 2) - a.Size
	if delta < 0 {
		return Atlas{}, fmt.Errorf("Cannot shrink an atlas")
	} else if delta == 0 {
		return Atlas{}, nil
	}

	// Create a new atlas
	newAtlas := NewAtlas(newRadius, a.ChunkSize)

	// Copy old chunks into new chunks
	for index, chunk := range a.Chunks {
		// Calculate the new chunk location and copy over the data
		newAtlas.Chunks[newAtlas.ToChunk(a.ChunkOrigin(index))] = chunk
	}

	// Return the new atlas
	return newAtlas, nil
}
