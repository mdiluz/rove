package atlas

import (
	"log"
	"math/rand"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
)

// chunk represents a fixed square grid of tiles
type chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles []byte `json:"tiles"`

	// Objects represents the objects within the chunk
	// only one possible object per tile for now
	Objects map[int]Object `json:"objects"`
}

// chunkBasedAtlas represents a grid of Chunks
type chunkBasedAtlas struct {
	// Chunks represents all chunks in the world
	// This is intentionally not a 2D array so it can be expanded in all directions
	Chunks []chunk `json:"chunks"`

	// LowerBound is the origin of the bottom left corner of the current chunks in world space (current chunks cover >= this value)
	LowerBound maths.Vector `json:"lowerBound"`

	// UpperBound is the top left corner of the current chunks (curent chunks cover < this value)
	UpperBound maths.Vector `json:"upperBound"`

	// ChunkSize is the x/y dimensions of each square chunk
	ChunkSize int `json:"chunksize"`

	// worldGen is the internal world generator
	worldGen WorldGen
}

const (
	noiseSeed = 1024
)

// NewChunkAtlas creates a new empty atlas
func NewChunkAtlas(chunkSize int) Atlas {
	// Start up with one chunk
	a := chunkBasedAtlas{
		ChunkSize:  chunkSize,
		Chunks:     make([]chunk, 1),
		LowerBound: maths.Vector{X: 0, Y: 0},
		UpperBound: maths.Vector{X: chunkSize, Y: chunkSize},
		worldGen:   NewNoiseWorldGen(noiseSeed),
	}
	// Initialise the first chunk
	a.populate(0)
	return &a
}

// SetTile sets an individual tile's kind
func (a *chunkBasedAtlas) SetTile(v maths.Vector, tile roveapi.Tile) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	a.setTile(c, local, byte(tile))
}

// SetObject sets the object on a tile
func (a *chunkBasedAtlas) SetObject(v maths.Vector, obj Object) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	a.setObject(c, local, obj)
}

// QueryPosition will return information for a specific position
func (a *chunkBasedAtlas) QueryPosition(v maths.Vector) (roveapi.Tile, Object) {
	c := a.worldSpaceToChunkWithGrow(v)
	local := a.worldSpaceToChunkLocal(v)
	a.populate(c)
	chunk := a.Chunks[c]
	i := a.chunkTileIndex(local)
	return roveapi.Tile(chunk.Tiles[i]), chunk.Objects[i]
}

// chunkTileID returns the tile index within a chunk
func (a *chunkBasedAtlas) chunkTileIndex(local maths.Vector) int {
	return local.X + local.Y*a.ChunkSize
}

// populate will fill a chunk with data
func (a *chunkBasedAtlas) populate(chunk int) {
	c := a.Chunks[chunk]
	if c.Tiles != nil {
		return
	}

	c.Tiles = make([]byte, a.ChunkSize*a.ChunkSize)
	c.Objects = make(map[int]Object)

	origin := a.chunkOriginInWorldSpace(chunk)
	for i := 0; i < a.ChunkSize; i++ {
		for j := 0; j < a.ChunkSize; j++ {
			loc := maths.Vector{X: origin.X + i, Y: origin.Y + j}

			// Set the tile
			c.Tiles[j*a.ChunkSize+i] = byte(a.worldGen.GetTile(loc))

			// Set the object
			obj := a.worldGen.GetObject(loc)
			if obj.Type != roveapi.Object_ObjectUnknown {
				c.Objects[j*a.ChunkSize+i] = obj
			}
		}
	}

	// Set up any objects
	for i := 0; i < len(c.Tiles); i++ {
		if rand.Intn(16) == 0 {
			c.Objects[i] = Object{Type: roveapi.Object_RockLarge}
		} else if rand.Intn(32) == 0 {
			c.Objects[i] = Object{Type: roveapi.Object_RockSmall}
		}
	}

	a.Chunks[chunk] = c
}

// setTile sets a tile in a specific chunk
func (a *chunkBasedAtlas) setTile(chunk int, local maths.Vector, tile byte) {
	a.populate(chunk)
	c := a.Chunks[chunk]
	c.Tiles[a.chunkTileIndex(local)] = tile
	a.Chunks[chunk] = c
}

// setObject sets an object in a specific chunk
func (a *chunkBasedAtlas) setObject(chunk int, local maths.Vector, object Object) {
	a.populate(chunk)

	c := a.Chunks[chunk]
	i := a.chunkTileIndex(local)
	if object.Type != roveapi.Object_ObjectUnknown {
		c.Objects[i] = object
	} else {
		delete(c.Objects, i)
	}
	a.Chunks[chunk] = c
}

// worldSpaceToChunkLocal gets a chunk local coordinate for a tile
func (a *chunkBasedAtlas) worldSpaceToChunkLocal(v maths.Vector) maths.Vector {
	return maths.Vector{X: maths.Pmod(v.X, a.ChunkSize), Y: maths.Pmod(v.Y, a.ChunkSize)}
}

// worldSpaceToChunkID gets the current chunk ID for a position in the world
func (a *chunkBasedAtlas) worldSpaceToChunkIndex(v maths.Vector) int {
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
func (a *chunkBasedAtlas) chunkOriginInWorldSpace(chunk int) maths.Vector {
	// Calculate the width
	width := a.UpperBound.X - a.LowerBound.X
	widthInChunks := width / a.ChunkSize

	// Reverse the along the corridor and up the stairs
	v := maths.Vector{
		X: chunk % widthInChunks,
		Y: chunk / widthInChunks,
	}
	// Multiply up to world scale
	v = v.Multiplied(a.ChunkSize)
	// Shift by the lower bound
	return v.Added(a.LowerBound)
}

// getNewBounds gets new lower and upper bounds for the world space given a vector
func (a *chunkBasedAtlas) getNewBounds(v maths.Vector) (lower maths.Vector, upper maths.Vector) {
	lower = maths.Min2(v, a.LowerBound)
	upper = maths.Max2(v.Added(maths.Vector{X: 1, Y: 1}), a.UpperBound)

	lower = maths.Vector{
		X: maths.RoundDown(lower.X, a.ChunkSize),
		Y: maths.RoundDown(lower.Y, a.ChunkSize),
	}
	upper = maths.Vector{
		X: maths.RoundUp(upper.X, a.ChunkSize),
		Y: maths.RoundUp(upper.Y, a.ChunkSize),
	}
	return
}

// worldSpaceToTrunkWithGrow will expand the current atlas for a given world space position if needed
func (a *chunkBasedAtlas) worldSpaceToChunkWithGrow(v maths.Vector) int {
	// If we're within bounds, just return the current chunk
	if v.X >= a.LowerBound.X && v.Y >= a.LowerBound.Y && v.X < a.UpperBound.X && v.Y < a.UpperBound.Y {
		return a.worldSpaceToChunkIndex(v)
	}

	// Calculate the new bounds
	lower, upper := a.getNewBounds(v)
	size := upper.Added(lower.Negated())
	size = size.Divided(a.ChunkSize)

	// Create the new empty atlas
	newAtlas := chunkBasedAtlas{
		ChunkSize:  a.ChunkSize,
		LowerBound: lower,
		UpperBound: upper,
		Chunks:     make([]chunk, size.X*size.Y),
		worldGen:   a.worldGen,
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
