package atlas

import (
	"fmt"
	"testing"

	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestAtlas_NewAtlas(t *testing.T) {
	a := NewChunkAtlas(1).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ChunkSize)
	assert.Equal(t, 1, len(a.Chunks)) // Should start empty
}

func TestAtlas_toChunk(t *testing.T) {
	a := NewChunkAtlas(1).(*ChunkBasedAtlas)
	assert.NotNil(t, a)

	// Get a tile to spawn the chunks
	a.QueryPosition(vector.Vector{X: -1, Y: -1})
	a.QueryPosition(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	//  2 | 3
	//  -----
	//  0 | 1
	chunkID := a.worldSpaceToChunkIndex(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, 3, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 0, Y: -1})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -1, Y: 0})
	assert.Equal(t, 2, chunkID)

	a = NewChunkAtlas(2).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.QueryPosition(vector.Vector{X: -2, Y: -2})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(vector.Vector{X: 1, Y: 1})
	assert.Equal(t, 2*2, len(a.Chunks))
	// Chunks should look like:
	// 2 | 3
	// -----
	// 0 | 1
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: 1})
	assert.Equal(t, 3, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -2, Y: -2})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -2, Y: 1})
	assert.Equal(t, 2, chunkID)

	a = NewChunkAtlas(2).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(vector.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(vector.Vector{X: -3, Y: -3})
	assert.Equal(t, 4*4, len(a.Chunks))

	// Chunks should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: 3})
	assert.Equal(t, 14, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: -3})
	assert.Equal(t, 2, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, 5, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: -2, Y: 2})
	assert.Equal(t, 13, chunkID)

	a = NewChunkAtlas(3).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(vector.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	// || 2| 3
	// -------
	// || 0| 1
	// =======
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: 1})
	assert.Equal(t, 0, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 3, Y: 1})
	assert.Equal(t, 1, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 1, Y: 4})
	assert.Equal(t, 2, chunkID)
	chunkID = a.worldSpaceToChunkIndex(vector.Vector{X: 5, Y: 5})
	assert.Equal(t, 3, chunkID)
}

func TestAtlas_toWorld(t *testing.T) {
	a := NewChunkAtlas(1).(*ChunkBasedAtlas)
	assert.NotNil(t, a)

	// Get a tile to spawn some chunks
	a.QueryPosition(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	//  2 | 3
	//  -----
	//  0 | 1
	assert.Equal(t, vector.Vector{X: -1, Y: -1}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, vector.Vector{X: 0, Y: -1}, a.chunkOriginInWorldSpace(1))

	a = NewChunkAtlas(2).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn the chunks
	a.QueryPosition(vector.Vector{X: -2, Y: -2})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(vector.Vector{X: 1, Y: 1})
	assert.Equal(t, 2*2, len(a.Chunks))
	// Chunks should look like:
	// 2 | 3
	// -----
	// 0 | 1
	assert.Equal(t, vector.Vector{X: -2, Y: -2}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, vector.Vector{X: -2, Y: 0}, a.chunkOriginInWorldSpace(2))

	a = NewChunkAtlas(2).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(vector.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))
	a.QueryPosition(vector.Vector{X: -3, Y: -3})
	assert.Equal(t, 4*4, len(a.Chunks))

	// Chunks should look like:
	//  12| 13|| 14| 15
	// ----------------
	//  8 | 9 || 10| 11
	// ================
	//  4 | 5 || 6 | 7
	// ----------------
	//  0 | 1 || 2 | 3
	assert.Equal(t, vector.Vector{X: -4, Y: -4}, a.chunkOriginInWorldSpace(0))
	assert.Equal(t, vector.Vector{X: 2, Y: -2}, a.chunkOriginInWorldSpace(7))

	a = NewChunkAtlas(3).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	// Get a tile to spawn a 4x4 grid of chunks
	a.QueryPosition(vector.Vector{X: 3, Y: 3})
	assert.Equal(t, 2*2, len(a.Chunks))

	// Chunks should look like:
	// || 2| 3
	// -------
	// || 0| 1
	// =======
	assert.Equal(t, vector.Vector{X: 0, Y: 0}, a.chunkOriginInWorldSpace(0))
}

func TestAtlas_GetSetTile(t *testing.T) {
	a := NewChunkAtlas(10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	a.SetTile(vector.Vector{X: 0, Y: 0}, 1)
	tile, _ := a.QueryPosition(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	// Set another tile to 1 and test it
	a.SetTile(vector.Vector{X: 5, Y: -2}, 2)
	tile, _ = a.QueryPosition(vector.Vector{X: 5, Y: -2})
	assert.Equal(t, byte(2), tile)
}

func TestAtlas_GetSetObject(t *testing.T) {
	a := NewChunkAtlas(10)
	assert.NotNil(t, a)

	// Set the origin tile to 1 and test it
	a.SetObject(vector.Vector{X: 0, Y: 0}, objects.Object{Type: objects.LargeRock})
	_, obj := a.QueryPosition(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, objects.Object{Type: objects.LargeRock}, obj)

	// Set another tile to 1 and test it
	a.SetObject(vector.Vector{X: 5, Y: -2}, objects.Object{Type: objects.SmallRock})
	_, obj = a.QueryPosition(vector.Vector{X: 5, Y: -2})
	assert.Equal(t, objects.Object{Type: objects.SmallRock}, obj)
}

func TestAtlas_Grown(t *testing.T) {
	// Start with a small example
	a := NewChunkAtlas(2).(*ChunkBasedAtlas)
	assert.NotNil(t, a)
	assert.Equal(t, 1, len(a.Chunks))

	// Set a few tiles to values
	a.SetTile(vector.Vector{X: 0, Y: 0}, 1)
	a.SetTile(vector.Vector{X: -1, Y: -1}, 2)
	a.SetTile(vector.Vector{X: 1, Y: -2}, 3)

	// Check tile values
	tile, _ := a.QueryPosition(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile, _ = a.QueryPosition(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile, _ = a.QueryPosition(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)

	tile, _ = a.QueryPosition(vector.Vector{X: 0, Y: 0})
	assert.Equal(t, byte(1), tile)

	tile, _ = a.QueryPosition(vector.Vector{X: -1, Y: -1})
	assert.Equal(t, byte(2), tile)

	tile, _ = a.QueryPosition(vector.Vector{X: 1, Y: -2})
	assert.Equal(t, byte(3), tile)
}

func TestAtlas_GetSetCorrect(t *testing.T) {
	// Big stress test to ensure we do actually properly expand for all reasonable values
	for i := 1; i <= 4; i++ {

		for x := -i * 2; x < i*2; x++ {
			for y := -i * 2; y < i*2; y++ {
				a := NewChunkAtlas(i).(*ChunkBasedAtlas)
				assert.NotNil(t, a)
				assert.Equal(t, 1, len(a.Chunks))

				pos := vector.Vector{X: x, Y: y}
				a.SetTile(pos, TileRock)
				a.SetObject(pos, objects.Object{Type: objects.LargeRock})
				tile, obj := a.QueryPosition(pos)

				assert.Equal(t, TileRock, Tile(tile))
				assert.Equal(t, objects.Object{Type: objects.LargeRock}, obj)

			}
		}
	}
}

func TestAtlas_WorldGen(t *testing.T) {
	a := NewChunkAtlas(8)
	// Spawn a large world
	_, _ = a.QueryPosition(vector.Vector{X: 20, Y: 20})

	// Print out the world for manual evaluation
	num := 20
	for j := num - 1; j >= 0; j-- {
		for i := 0; i < num; i++ {
			t, o := a.QueryPosition(vector.Vector{X: i, Y: j})
			if o.Type != objects.None {
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
