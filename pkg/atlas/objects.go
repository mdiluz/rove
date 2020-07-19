package atlas

import (
	"github.com/mdiluz/rove/proto/roveapi"
)

// Object represents an object in the world
type Object struct {
	// The type of the object
	Type roveapi.Object `json:"type"`
}

// IsBlocking checks if an object is a blocking object
func (o *Object) IsBlocking() bool {
	var blocking = [...]roveapi.Object{
		roveapi.Object_RoverLive,
		roveapi.Object_RockLarge,
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
	var stashable = [...]roveapi.Object{
		roveapi.Object_RockSmall,
	}

	for _, t := range stashable {
		if o.Type == t {
			return true
		}
	}
	return false
}
