package atlas

// Glyph represents the text representation of something in the game
type Glyph byte

const (
	// GlyphNone is a keyword for nothing
	GlyphNone = Glyph(0)

	// GlyphRock is solid rock ground
	GlyphRock = Glyph('-')

	// GlyphGravel is loose rocks
	GlyphGravel = Glyph(':')

	// GlyphSand is sand
	GlyphSand = Glyph('~')

	// GlyphRover represents a live rover
	GlyphRover = Glyph('R')

	// GlyphSmallRock is a small stashable rock
	GlyphSmallRock = Glyph('o')

	// GlyphLargeRock is a large blocking rock
	GlyphLargeRock = Glyph('O')
)
