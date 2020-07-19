package atlas

import (
	"log"

	"github.com/mdiluz/rove/pkg/maths"
)

// Tile describes the type of terrain
type Tile byte

const (
	// TileNone is a keyword for nothing
	TileNone = iota

	// TileRock is solid rock ground
	TileRock

	// TileGravel is loose rocks
	TileGravel

	// TileSand is sand
	TileSand
)

// Glyph returns the glyph for this tile type
func (t Tile) Glyph() Glyph {
	switch t {
	case TileNone:
		return GlyphNone
	case TileRock:
		return GlyphGroundRock
	case TileGravel:
		return GlyphGroundGravel
	case TileSand:
		return GlyphGroundSand
	}

	log.Fatalf("Unknown tile type: %c", t)
	return GlyphNone
}

// Atlas represents a 2D world atlas of tiles and objects
type Atlas interface {
	// SetTile sets a location on the Atlas to a type of tile
	SetTile(v maths.Vector, tile Tile)

	// SetObject will set a location on the Atlas to contain an object
	SetObject(v maths.Vector, obj Object)

	// QueryPosition queries a position on the atlas
	QueryPosition(v maths.Vector) (byte, Object)
}
