package atlas

import (
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

	// CurrentSize is the current width/height of the atlas in chunks
	CurrentSize vector.Vector `json:"currentSize"`

	// WorldOrigin represents the location of the [0,0] world space point in terms of the allotted current chunks
	WorldOrigin vector.Vector `json:"worldOrigin"`

	// ChunkSize is the x/y dimensions of each square chunk
	ChunkSize int `json:"chunksize"`
}

// NewAtlas creates a new empty atlas
func NewAtlas(chunkSize int) Atlas {
	// Start up with one chunk
	a := Atlas{
		ChunkSize:   chunkSize,
		Chunks:      make([]Chunk, 1),
		CurrentSize: vector.Vector{X: 1, Y: 1},
		WorldOrigin: vector.Vector{X: 0, Y: 0},
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

// setTileAndObject sets both tile and object information for location in chunk
func (a *Atlas) setTileAndObject(chunk int, local vector.Vector, tile byte, object objects.Object) {
	c := a.Chunks[chunk]
	if c.Tiles == nil {
		c.populate(a.ChunkSize)
	}

	i := a.chunkTileIndex(local)
	c.Tiles[i] = tile
	c.Objects[i] = object
	a.Chunks[chunk] = c
}

// worldSpaceToChunkLocal gets a chunk local coordinate for a tile
func (a *Atlas) worldSpaceToChunkLocal(v vector.Vector) vector.Vector {
	return vector.Vector{X: maths.Pmod(v.X, a.ChunkSize), Y: maths.Pmod(v.Y, a.ChunkSize)}
}

// worldSpaceToChunk gets the current chunk ID for a position in the world
func (a *Atlas) worldSpaceToChunk(v vector.Vector) int {
	// First convert to chunk space
	chunkSpace := a.worldSpaceToChunkSpace(v)

	// Then return the ID
	return a.chunkSpaceToChunk(chunkSpace)
}

// worldSpaceToChunkSpace converts from world space to chunk space
func (a *Atlas) worldSpaceToChunkSpace(v vector.Vector) vector.Vector {
	// Remove the chunk local part
	chunkOrigin := v.Added(a.worldSpaceToChunkLocal(v).Negated())
	// Convert to chunk space coordinate
	chunkSpaceOrigin := chunkOrigin.Divided(a.ChunkSize)
	// Shift it by our current chunk origin
	chunkIndexOrigin := chunkSpaceOrigin.Added(a.WorldOrigin)

	return chunkIndexOrigin
}

// chunkSpaceToWorldSpace vonverts from chunk space to world space
func (a *Atlas) chunkSpaceToWorldSpace(v vector.Vector) vector.Vector {

	// Shift it by the current chunk origin
	shifted := v.Added(a.WorldOrigin.Negated())

	// Multiply out by chunk size
	return shifted.Multiplied(a.ChunkSize)
}

// chunkOriginInChunkSpace Gets the chunk origin in chunk space
func (a *Atlas) chunkOriginInChunkSpace(chunk int) vector.Vector {
	// convert the chunk to chunk space
	chunkOrigin := a.chunkToChunkSpace(chunk)

	// Shift it by the current chunk origin
	return chunkOrigin.Added(a.WorldOrigin.Negated())
}

// chunkOriginInWorldSpace gets the chunk origin for a given chunk index
func (a *Atlas) chunkOriginInWorldSpace(chunk int) vector.Vector {
	// convert the chunk to chunk space
	chunkSpace := a.chunkToChunkSpace(chunk)

	// Convert to world space
	return a.chunkSpaceToWorldSpace(chunkSpace)
}

// chunkSpaceToChunk converts from chunk space to the chunk
func (a *Atlas) chunkSpaceToChunk(v vector.Vector) int {
	// Along the coridor and up the stair
	return (v.Y * a.CurrentSize.X) + v.X
}

// chunkToChunkSpace returns the chunk space coord for the chunk
func (a *Atlas) chunkToChunkSpace(chunk int) vector.Vector {
	return vector.Vector{
		X: maths.Pmod(chunk, a.CurrentSize.Y),
		Y: (chunk / a.CurrentSize.X),
	}
}

func (a *Atlas) getExtents() (min vector.Vector, max vector.Vector) {
	min = a.WorldOrigin.Negated()
	max = min.Added(a.CurrentSize)
	return
}

// worldSpaceToTrunkWithGrow will expand the current atlas for a given world space position if needed
func (a *Atlas) worldSpaceToChunkWithGrow(v vector.Vector) int {
	min, max := a.getExtents()

	// Divide by the chunk size to bring into chunk space
	v = v.Divided(a.ChunkSize)

	// Check we're within the current extents and bail early
	if v.X >= min.X && v.Y >= min.Y && v.X < max.X && v.Y < max.Y {
		return a.worldSpaceToChunk(v)
	}

	// Calculate the new origin and the new size
	origin := min
	size := a.CurrentSize

	// If we need to shift the origin back
	originDiff := origin.Added(v.Negated())
	if originDiff.X > 0 {
		origin.X -= originDiff.X
		size.X += originDiff.X
	}
	if originDiff.Y > 0 {
		origin.Y -= originDiff.Y
		size.Y += originDiff.Y
	}

	// If we need to expand the size
	maxDiff := v.Added(max.Negated())
	if maxDiff.X > 0 {
		size.X += maxDiff.X
	}
	if maxDiff.Y > 0 {
		size.Y += maxDiff.Y
	}

	// Set up the new size and origin
	newAtlas := Atlas{
		ChunkSize:   a.ChunkSize,
		WorldOrigin: origin.Negated(),
		CurrentSize: size,
		Chunks:      make([]Chunk, size.X*size.Y),
	}

	// Copy all old chunks into the new atlas
	for chunk, chunkData := range a.Chunks {
		// Calculate the new chunk location and copy over the data
		newChunk := newAtlas.worldSpaceToChunk(a.chunkOriginInWorldSpace(chunk))
		// Copy over the old chunk to the new atlas
		newAtlas.Chunks[newChunk] = chunkData
	}

	// Copy the new atlas data into this one
	*a = newAtlas

	return a.worldSpaceToChunk(v)
}
