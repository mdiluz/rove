package atlas

import (
	"log"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
)

// TileGlyph returns the glyph for this tile type
func TileGlyph(t roveapi.Tile) Glyph {
	switch t {
	case roveapi.Tile_TileNone:
		return GlyphNone
	case roveapi.Tile_Rock:
		return GlyphGroundRock
	case roveapi.Tile_Gravel:
		return GlyphGroundGravel
	case roveapi.Tile_Sand:
		return GlyphGroundSand
	}

	log.Fatalf("Unknown tile type: %c", t)
	return GlyphNone
}

// Atlas represents a 2D world atlas of tiles and objects
type Atlas interface {
	// SetTile sets a location on the Atlas to a type of tile
	SetTile(v maths.Vector, tile roveapi.Tile)

	// SetObject will set a location on the Atlas to contain an object
	SetObject(v maths.Vector, obj Object)

	// QueryPosition queries a position on the atlas
	QueryPosition(v maths.Vector) (roveapi.Tile, Object)
}
