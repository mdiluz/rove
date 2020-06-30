package objects

const (
	// Empty represents an non-existant object
	Empty = byte(' ')

	// Rover represents a live rover
	Rover = byte('R')

	// SmallRock is a small stashable rock
	SmallRock = byte('o')

	// LargeRock is a large blocking rock
	LargeRock = byte('O')
)

// IsBlocking checks if an object is a blocking object
func IsBlocking(object byte) bool {
	var blocking = [...]byte{
		Rover,
		LargeRock,
	}

	for _, t := range blocking {
		if object == t {
			return true
		}
	}
	return false
}

// IsStashable checks if an object is stashable
func IsStashable(object byte) bool {
	var stashable = [...]byte{
		SmallRock,
	}

	for _, t := range stashable {
		if object == t {
			return true
		}
	}
	return false
}
