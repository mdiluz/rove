package atlas

import (
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
)

// Atlas represents a 2D world atlas of tiles and objects
type Atlas interface {
	// SetTile sets a location on the Atlas to a type of tile
	SetTile(v maths.Vector, tile roveapi.Tile)

	// SetObject will set a location on the Atlas to contain an object
	SetObject(v maths.Vector, obj Object)

	// QueryPosition queries a position on the atlas
	QueryPosition(v maths.Vector) (roveapi.Tile, Object)
}
