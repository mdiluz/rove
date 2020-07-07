package atlas

import (
	"log"
	"math/rand"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
)

// Tile describes the type of terrain
type Tile byte

const (
	// TileNone is a keyword for nothing
	TileNone = Tile(0)

	// TileRock is solid rock ground
	TileRock = Tile('.')

	// TileSand is sand
	TileSand = Tile(',')
)

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []byte `json:"tiles"`

	// Objects represents the objects within the chunk
	// only one possible object per tile for now
	Objects map[int]objects.Object `json:"objects"`
}

// Atlas represents a grid of Chunks
type Atlas struct {
	// Chunks represents all chunks in the world
	// This is intentionally not a 2D array so it can be expanded in all directions
	Chunks []Chunk `json:"chunks"`

	// LowerBound is the origin of the bottom left corner of the current chunks in world space (current chunks cover >= this value)
	LowerBound vector.Vector `json:"lowerBound"`

	// UpperBound is the top left corner of the current chunks (curent chunks cover < this value)
	UpperBound vector.Vector `json:"upperBound"`

	// ChunkSize is the x/y dimensions of each square chunk
	ChunkSize int `json:"chunksize"`
}

// NewAtlas creates a new empty atlas
func NewAtlas(chunkSize int) Atlas {
	// Start up with one chunk
	a := Atlas{
		ChunkSize:  chunkSize,
		Chunks:     make([]Chunk, 1),
		LowerBound: vector.Vector{X: 0, Y: 0},
		UpperBound: vector.Vector{X: chunkSize, Y: chunkSize},
	}
	// Initialise the first chunk
	a.Chunks[0].populate(chunkSize)
	return a
}

// SetTile sets an individual tile's kind
func (a *Atlas) SetTile(v vector.Vector, tile Tile) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	a.setTile(c, local, byte(tile))
}

// SetObject sets the object on a tile
func (a *Atlas) SetObject(v vector.Vector, obj objects.Object) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	a.setObject(c, local, obj)
}

// QueryPosition will return information for a specific position
func (a *Atlas) QueryPosition(v vector.Vector) (byte, objects.Object) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	chunk := a.Chunks[c]
	if chunk.Tiles == nil {
		chunk.populate(a.ChunkSize)
	}
	i := a.chunkTileIndex(local)
	return chunk.Tiles[i], chunk.Objects[i]
}

// chunkTileID returns the tile index within a chunk
func (a *Atlas) chunkTileIndex(local vector.Vector) int {
	return local.X + local.Y*a.ChunkSize
}

// populate will fill a chunk with data
func (c *Chunk) populate(size int) {
	c.Tiles = make([]byte, size*size)
	c.Objects = make(map[int]objects.Object)

	// Set up the tiles
	for i := 0; i < len(c.Tiles); i++ {
		if rand.Intn(3) == 0 {
			c.Tiles[i] = byte(TileRock)
		} else {
			c.Tiles[i] = byte(TileSand)
		}
	}

	// Set up any objects
	for i := 0; i < len(c.Tiles); i++ {
		if rand.Intn(16) == 0 {
			c.Objects[i] = objects.Object{Type: objects.LargeRock}
		} else if rand.Intn(32) == 0 {
			c.Objects[i] = objects.Object{Type: objects.SmallRock}
		}
	}
}

// setTile sets a tile in a specific chunk
func (a *Atlas) setTile(chunk int, local vector.Vector, tile byte) {
	c := a.Chunks[chunk]
	if c.Tiles == nil {
		c.populate(a.ChunkSize)
	}

	c.Tiles[a.chunkTileIndex(local)] = tile
	a.Chunks[chunk] = c
}

// setObject sets an object in a specific chunk
func (a *Atlas) setObject(chunk int, local vector.Vector, object objects.Object) {
	c := a.Chunks[chunk]
	if c.Tiles == nil {
		c.populate(a.ChunkSize)
	}

	i := a.chunkTileIndex(local)
	if object.Type != objects.None {
		c.Objects[i] = object
	} else {
		delete(c.Objects, i)
	}
	a.Chunks[chunk] = c
}

// worldSpaceToChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) worldSpaceToChunkLocal(v vector.Vector) vector.Vector {
	return vector.Vector{X: maths.Pmod(v.X, a.ChunkSize), Y: maths.Pmod(v.Y, a.ChunkSize)}
}

// worldSpaceToChunkID gets the current chunk ID for a position in the world
func (a *Atlas) worldSpaceToChunkIndex(v vector.Vector) int {
	// Shift the vector by our current min
	v = v.Added(a.LowerBound.Negated())

	// Divide by the current size and floor, to get chunk-scaled vector from the lower bound
	v = v.DividedFloor(a.ChunkSize)

	// Calculate the width
	width := a.UpperBound.X - a.LowerBound.X
	widthInChunks := width / a.ChunkSize

	// Along the corridor and up the stairs
	return (v.Y * widthInChunks) + v.X
}

// chunkOriginInWorldSpace returns the origin of the chunk in world space
func (a *Atlas) chunkOriginInWorldSpace(chunk int) vector.Vector {
	// Calculate the width
	width := a.UpperBound.X - a.LowerBound.X
	widthInChunks := width / a.ChunkSize

	// Reverse the along the corridor and up the stairs
	v := vector.Vector{
		X: chunk % widthInChunks,
		Y: chunk / widthInChunks,
	}
	// Multiply up to world scale
	v = v.Multiplied(a.ChunkSize)
	// Shift by the lower bound
	return v.Added(a.LowerBound)
}

// getNewBounds gets new lower and upper bounds for the world space given a vector
func (a *Atlas) getNewBounds(v vector.Vector) (lower vector.Vector, upper vector.Vector) {
	lower = vector.Min(v, a.LowerBound)
	upper = vector.Max(v.Added(vector.Vector{X: 1, Y: 1}), a.UpperBound)

	lower = vector.Vector{
		X: maths.RoundDown(lower.X, a.ChunkSize),
		Y: maths.RoundDown(lower.Y, a.ChunkSize),
	}
	upper = vector.Vector{
		X: maths.RoundUp(upper.X, a.ChunkSize),
		Y: maths.RoundUp(upper.Y, a.ChunkSize),
	}
	return
}

// worldSpaceToTrunkWithGrow will expand the current atlas for a given world space position if needed
func (a *Atlas) worldSpaceToChunkWithGrow(v vector.Vector) int {
	// If we're within bounds, just return the current chunk
	if v.X >= a.LowerBound.X && v.Y >= a.LowerBound.Y && v.X < a.UpperBound.X && v.Y < a.UpperBound.Y {
		return a.worldSpaceToChunkIndex(v)
	}

	// Calculate the new bounds
	lower, upper := a.getNewBounds(v)
	size := upper.Added(lower.Negated())
	size = size.Divided(a.ChunkSize)

	// Create the new empty atlas
	newAtlas := Atlas{
		ChunkSize:  a.ChunkSize,
		LowerBound: lower,
		UpperBound: upper,
		Chunks:     make([]Chunk, size.X*size.Y),
	}

	// Log that we're resizing
	log.Printf("Re-allocating world, old: %+v,%+v new: %+v,%+v\n", a.LowerBound, a.UpperBound, newAtlas.LowerBound, newAtlas.UpperBound)

	// Copy all old chunks into the new atlas
	for chunk, chunkData := range a.Chunks {

		// Calculate the chunk ID in the new atlas
		origin := a.chunkOriginInWorldSpace(chunk)
		newChunk := newAtlas.worldSpaceToChunkIndex(origin)

		// Copy over the old chunk to the new atlas
		newAtlas.Chunks[newChunk] = chunkData
	}

	// Overwrite the old atlas with this one
	*a = newAtlas

	return a.worldSpaceToChunkIndex(v)
}
