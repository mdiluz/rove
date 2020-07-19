package atlas

// Type represents an object type
type Type byte

// Types of objects
const (
	// ObjectNone represents no object at all
	ObjectNone = Type(GlyphNone)

	// ObjectRover represents a live rover
	ObjectRover = Type(GlyphRover)

	// ObjectSmallRock is a small stashable rock
	ObjectSmallRock = Type(GlyphSmallRock)

	// ObjectLargeRock is a large blocking rock
	ObjectLargeRock = Type('O')
)

// Object represents an object in the world
type Object struct {
	Type Type `json:"type"`
}

// IsBlocking checks if an object is a blocking object
func (o *Object) IsBlocking() bool {
	var blocking = [...]Type{
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
	var stashable = [...]Type{
		ObjectSmallRock,
	}

	for _, t := range stashable {
		if o.Type == t {
			return true
		}
	}
	return false
}
