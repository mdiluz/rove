package atlas

import "log"

// ObjectType represents an object type
type ObjectType byte

// Types of objects
const (
	// ObjectNone represents no object at all
	ObjectNone = iota

	// ObjectRover represents a live rover
	ObjectRover

	// ObjectSmallRock is a small stashable rock
	ObjectSmallRock

	// ObjectLargeRock is a large blocking rock
	ObjectLargeRock
)

// Glyph returns the glyph for this object type
func (o ObjectType) Glyph() Glyph {
	switch o {
	case ObjectNone:
		return GlyphNone
	case ObjectRover:
		return GlyphRover
	case ObjectSmallRock:
		return GlyphSmallRock
	case ObjectLargeRock:
		return GlyphLargeRock
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
		ObjectRover,
		ObjectLargeRock,
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
		ObjectSmallRock,
	}

	for _, t := range stashable {
		if o.Type == t {
			return true
		}
	}
	return false
}
