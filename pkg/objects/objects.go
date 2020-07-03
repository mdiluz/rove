package objects

// Type represents an object type
type Type byte

// Types of objects
const (
	// None represents no object at all
	None = Type(0)

	// Rover represents a live rover
	Rover = Type('R')

	// SmallRock is a small stashable rock
	SmallRock = Type('o')

	// LargeRock is a large blocking rock
	LargeRock = Type('O')
)

// Object represents an object in the world
type Object struct {
	Type Type `json:"type"`
}

// IsBlocking checks if an object is a blocking object
func (o *Object) IsBlocking() bool {
	var blocking = [...]Type{
		Rover,
		LargeRock,
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
		SmallRock,
	}

	for _, t := range stashable {
		if o.Type == t {
			return true
		}
	}
	return false
}
