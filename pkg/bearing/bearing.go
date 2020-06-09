package bearing

import (
	"fmt"
	"strings"

	"github.com/mdiluz/rove/pkg/vector"
)

// Direction describes a compass direction
type Direction int

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

// DirectionString simply describes the strings associated with a direction
type DirectionString struct {
	Long  string
	Short string
}

// DirectionStrings is the set of strings for each direction
var DirectionStrings = []DirectionString{
	{"North", "N"},
	{"NorthEast", "NE"},
	{"East", "E"},
	{"SouthEast", "SE"},
	{"South", "S"},
	{"SouthWest", "SW"},
	{"West", "W"},
	{"NorthWest", "NW"},
}

// String converts a Direction to a String
func (d Direction) String() string {
	return DirectionStrings[d].Long
}

// ShortString converts a Direction to a short string version
func (d Direction) ShortString() string {
	return DirectionStrings[d].Short
}

// DirectionFromString gets the Direction from a string
func DirectionFromString(s string) (Direction, error) {
	for i, d := range DirectionStrings {
		if strings.ToLower(d.Long) == strings.ToLower(s) || strings.ToLower(d.Short) == strings.ToLower(s) {
			return Direction(i), nil
		}
	}
	return -1, fmt.Errorf("Unknown direction: %s", s)
}

var DirectionVectors = []vector.Vector{
	{X: 0, Y: 1},  // N
	{X: 1, Y: 1},  // NE
	{X: 1, Y: 0},  // E
	{X: 1, Y: -1}, // SE
	{X: 0, Y: -1}, // S
	{X: -1, Y: 1}, // SW
	{X: -1, Y: 0}, // W
	{X: -1, Y: 1}, // NW
}

// Vector converts a Direction to a Vector
func (d Direction) Vector() vector.Vector {
	return DirectionVectors[d]
}
