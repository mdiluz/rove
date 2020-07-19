package atlas

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
