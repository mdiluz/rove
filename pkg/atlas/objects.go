package atlas

import "log"

// ObjectType represents an object type
type ObjectType byte

// Types of objects
const (
	// ObjectNone represents no object at all
	ObjectNone = iota

	// ObjectRover represents a live rover
	ObjectRoverLive

	// ObjectSmallRock is a small stashable rock
	ObjectRockSmall

	// ObjectLargeRock is a large blocking rock
	ObjectRockLarge
)

// Glyph returns the glyph for this object type
func (o ObjectType) Glyph() Glyph {
	switch o {
	case ObjectNone:
		return GlyphNone
	case ObjectRoverLive:
		return GlyphRoverLive
	case ObjectRockSmall:
		return GlyphRockSmall
	case ObjectRockLarge:
		return GlyphRockLarge
	}

	log.Fatalf("Unknown object type: %c", o)
	return GlyphNone
}

// Object represents an object in the world
type Object struct {
	// The type of the object
	Type ObjectType `json:"type"`
}

// IsBlocking checks if an object is a blocking object
func (o *Object) IsBlocking() bool {
	var blocking = [...]ObjectType{
		ObjectRoverLive,
		ObjectRockLarge,
	}

	for _, t := range blocking {
		if o.Type == t {
			return true
		}
	}
	return false
}

// IsStashable checks if an object is stashable
func (o *Object) IsStashable() bool {
	var stashable = [...]ObjectType{
		ObjectRockSmall,
	}

	for _, t := range stashable {
		if o.Type == t {
			return true
		}
	}
	return false
}
