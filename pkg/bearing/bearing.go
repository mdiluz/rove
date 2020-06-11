package bearing

import (
	"fmt"
	"strings"

	"github.com/mdiluz/rove/pkg/vector"
)

// Bearing describes a compass direction
type Bearing int

const (
	North Bearing = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

// bearingString simply describes the strings associated with a direction
type bearingString struct {
	Long  string
	Short string
}

// bearingStrings is the set of strings for each direction
var bearingStrings = []bearingString{
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
func (d Bearing) String() string {
	return bearingStrings[d].Long
}

// ShortString converts a Direction to a short string version
func (d Bearing) ShortString() string {
	return bearingStrings[d].Short
}

// FromString gets the Direction from a string
func FromString(s string) (Bearing, error) {
	for i, d := range bearingStrings {
		if strings.ToLower(d.Long) == strings.ToLower(s) || strings.ToLower(d.Short) == strings.ToLower(s) {
			return Bearing(i), nil
		}
	}
	return -1, fmt.Errorf("unknown bearing: %s", s)
}

var bearingVectors = []vector.Vector{
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
func (d Bearing) Vector() vector.Vector {
	return bearingVectors[d]
}
