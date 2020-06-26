package objects

const (
	Empty     = byte(' ')
	Rover     = byte('R')
	SmallRock = byte('o')
	LargeRock = byte('O')
)

// Check if an object is a blocking object
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

// Check if an object is stashable
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
