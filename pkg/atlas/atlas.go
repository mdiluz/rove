package atlas

import (
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
)

// Tile describes the type of terrain
type Tile byte

const (
	// TileNone is a keyword for nothing
	TileNone = Tile(0)

	// TileRock is solid rock ground
	TileRock = Tile('-')

	// TileGravel is loose rocks
	TileGravel = Tile(':')

	// TileSand is sand
	TileSand = Tile('~')
)

// Atlas represents a 2D world atlas of tiles and objects
type Atlas interface {
	// SetTile sets a location on the Atlas to a type of tile
	SetTile(v vector.Vector, tile Tile)

	// SetObject will set a location on the Atlas to contain an object
	SetObject(v vector.Vector, obj objects.Object)

	// QueryPosition queries a position on the atlas
	QueryPosition(v vector.Vector) (byte, objects.Object)
}
