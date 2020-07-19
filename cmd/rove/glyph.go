package main

import (
	"log"

	"github.com/mdiluz/rove/proto/roveapi"
)

// Glyph represents the text representation of something in the game
type Glyph byte

const (
	// GlyphGroundRock is solid rock ground
	GlyphGroundRock = Glyph('-')

	// GlyphGroundGravel is loose rocks
	GlyphGroundGravel = Glyph(':')

	// GlyphGroundSand is sand
	GlyphGroundSand = Glyph('~')

	// GlyphRoverLive represents a live rover
	GlyphRoverLive = Glyph('R')

	// GlyphRockSmall is a small stashable rock
	GlyphRockSmall = Glyph('o')

	// GlyphRockLarge is a large blocking rock
	GlyphRockLarge = Glyph('O')
)

// TileGlyph returns the glyph for this tile type
func TileGlyph(t roveapi.Tile) Glyph {
	switch t {
	case roveapi.Tile_Rock:
		return GlyphGroundRock
	case roveapi.Tile_Gravel:
		return GlyphGroundGravel
	case roveapi.Tile_Sand:
		return GlyphGroundSand
	}

	log.Fatalf("Unknown tile type: %c", t)
	return 0
}

// ObjectGlyph returns the glyph for this object type
func ObjectGlyph(o roveapi.Object) Glyph {
	switch o {
	case roveapi.Object_RoverLive:
		return GlyphRoverLive
	case roveapi.Object_RockSmall:
		return GlyphRockSmall
	case roveapi.Object_RockLarge:
		return GlyphRockLarge
	}

	log.Fatalf("Unknown object type: %c", o)
	return 0
}
