package atlas

import (
	"fmt"
	"testing"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/stretchr/testify/assert"
)

func TestAtlas_NewAtlas(t *testing.T) {
	a := NewChunkAtlas(1).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ChunkSize)
	assert.Equal(t, 1, len(a.Chunks)) // Should start empty
}

func TestAtlas_toChunk(t *testing.T) {
	a := NewChunkAtlas(1).(*chunkBasedAtlas)
	assert.NotNil(t, a)

	// Get a tile to spawn the chunks
	a.QueryPosition(maths.Vector{X: -1, Y: -1})
	a.QueryPosition(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	//  2 | 3
	//  -----
	//  0 | 1
	chunkID := a.worldSpaceToChunkIndex(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, 3, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 0, Y: -1})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -1, Y: -1})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -1, Y: 0})
	assert.Equal(t, 2, chunkID)

	a = NewChunkAtlas(2).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.QueryPosition(maths.Vector{X: -2, Y: -2})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(maths.Vector{X: 1, Y: 1})
	assert.Equal(t, 2*2, len(a.Chunks))
	// Chunks should look like:
	// 2 | 3
	// -----
	// 0 | 1
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: 1})
	assert.Equal(t, 3, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: -2})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -2, Y: -2})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -2, Y: 1})
	assert.Equal(t, 2, chunkID)

	a = NewChunkAtlas(2).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(maths.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(maths.Vector{X: -3, Y: -3})
	assert.Equal(t, 4*4, len(a.Chunks))

	// Chunks should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: 3})
	assert.Equal(t, 14, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: -3})
	assert.Equal(t, 2, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -1, Y: -1})
	assert.Equal(t, 5, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: -2, Y: 2})
	assert.Equal(t, 13, chunkID)

	a = NewChunkAtlas(3).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(maths.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	// || 2| 3
	// -------
	// || 0| 1
	// =======
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: 1})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 3, Y: 1})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 1, Y: 4})
	assert.Equal(t, 2, chunkID)
	chunkID = a.worldSpaceToChunkIndex(maths.Vector{X: 5, Y: 5})
	assert.Equal(t, 3, chunkID)
}

func TestAtlas_toWorld(t *testing.T) {
	a := NewChunkAtlas(1).(*chunkBasedAtlas)
	assert.NotNil(t, a)

	// Get a tile to spawn some chunks
	a.QueryPosition(maths.Vector{X: -1, Y: -1})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	//  2 | 3
	//  -----
	//  0 | 1
	assert.Equal(t, maths.Vector{X: -1, Y: -1}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, maths.Vector{X: 0, Y: -1}, a.chunkOriginInWorldSpace(1))

	a = NewChunkAtlas(2).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.QueryPosition(maths.Vector{X: -2, Y: -2})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(maths.Vector{X: 1, Y: 1})
	assert.Equal(t, 2*2, len(a.Chunks))
	// Chunks should look like:
	// 2 | 3
	// -----
	// 0 | 1
	assert.Equal(t, maths.Vector{X: -2, Y: -2}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, maths.Vector{X: -2, Y: 0}, a.chunkOriginInWorldSpace(2))

	a = NewChunkAtlas(2).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(maths.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(maths.Vector{X: -3, Y: -3})
	assert.Equal(t, 4*4, len(a.Chunks))

	// Chunks should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	assert.Equal(t, maths.Vector{X: -4, Y: -4}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, maths.Vector{X: 2, Y: -2}, a.chunkOriginInWorldSpace(7))

	a = NewChunkAtlas(3).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(maths.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	// || 2| 3
	// -------
	// || 0| 1
	// =======
	assert.Equal(t, maths.Vector{X: 0, Y: 0}, a.chunkOriginInWorldSpace(0))
}

func TestAtlas_GetSetTile(t *testing.T) {
	a := NewChunkAtlas(10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	a.SetTile(maths.Vector{X: 0, Y: 0}, 1)
	tile, _ := a.QueryPosition(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	// Set another tile to 1 and test it
	a.SetTile(maths.Vector{X: 5, Y: -2}, 2)
	tile, _ = a.QueryPosition(maths.Vector{X: 5, Y: -2})
	assert.Equal(t, byte(2), tile)
}

func TestAtlas_GetSetObject(t *testing.T) {
	a := NewChunkAtlas(10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	a.SetObject(maths.Vector{X: 0, Y: 0}, Object{Type: ObjectRockLarge})
	_, obj := a.QueryPosition(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, Object{Type: ObjectRockLarge}, obj)

	// Set another tile to 1 and test it
	a.SetObject(maths.Vector{X: 5, Y: -2}, Object{Type: ObjectRockSmall})
	_, obj = a.QueryPosition(maths.Vector{X: 5, Y: -2})
	assert.Equal(t, Object{Type: ObjectRockSmall}, obj)
}

func TestAtlas_Grown(t *testing.T) {
	// Start with a small example
	a := NewChunkAtlas(2).(*chunkBasedAtlas)
	assert.NotNil(t, a)
	assert.Equal(t, 1, len(a.Chunks))

	// Set a few tiles to values
	a.SetTile(maths.Vector{X: 0, Y: 0}, 1)
	a.SetTile(maths.Vector{X: -1, Y: -1}, 2)
	a.SetTile(maths.Vector{X: 1, Y: -2}, 3)

	// Check tile values
	tile, _ := a.QueryPosition(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile, _ = a.QueryPosition(maths.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile, _ = a.QueryPosition(maths.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)

	tile, _ = a.QueryPosition(maths.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile, _ = a.QueryPosition(maths.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile, _ = a.QueryPosition(maths.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)
}

func TestAtlas_GetSetCorrect(t *testing.T) {
	// Big stress test to ensure we do actually properly expand for all reasonable values
	for i := 1; i <= 4; i++ {

		for x := -i * 2; x < i*2; x++ {
			for y := -i * 2; y < i*2; y++ {
				a := NewChunkAtlas(i).(*chunkBasedAtlas)
				assert.NotNil(t, a)
				assert.Equal(t, 1, len(a.Chunks))

				pos := maths.Vector{X: x, Y: y}
				a.SetTile(pos, TileRock)
				a.SetObject(pos, Object{Type: ObjectRockLarge})
				tile, obj := a.QueryPosition(pos)

				assert.Equal(t, TileRock, Tile(tile))
				assert.Equal(t, Object{Type: ObjectRockLarge}, obj)

			}
		}
	}
}

func TestAtlas_WorldGen(t *testing.T) {
	a := NewChunkAtlas(8)
	// Spawn a large world
	_, _ = a.QueryPosition(maths.Vector{X: 20, Y: 20})

	// Print out the world for manual evaluation
	num := 20
	for j := num - 1; j >= 0; j-- {
		for i := 0; i < num; i++ {
			t, o := a.QueryPosition(maths.Vector{X: i, Y: j})
			if o.Type != ObjectNone {
				fmt.Printf("%c", o.Type)
			} else if t != byte(TileNone) {
				fmt.Printf("%c", t)
			} else {
				fmt.Printf(" ")
			}

		}
		fmt.Print("\n")
	}
}
