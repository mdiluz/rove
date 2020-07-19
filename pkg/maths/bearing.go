package maths

import (
	"fmt"
	"strings"
)

// Bearing describes a compass direction
type Bearing int

const (
	// North describes a 0,1 vector
	North Bearing = iota
	// NorthEast describes a 1,1 vector
	NorthEast
	// East describes a 1,0 vector
	East
	// SouthEast describes a 1,-1 vector
	SouthEast
	// South describes a 0,-1 vector
	South
	// SouthWest describes a -1,-1 vector
	SouthWest
	// West describes a -1,0 vector
	West
	// NorthWest describes a -1,1 vector
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

// BearingFromString gets the Direction from a string
func BearingFromString(s string) (Bearing, error) {
	for i, d := range bearingStrings {
		if strings.EqualFold(d.Long, s) || strings.EqualFold(d.Short, s) {
			return Bearing(i), nil
		}
	}
	return -1, fmt.Errorf("unknown bearing: %s", s)
}

var bearingVectors = []Vector{
	{X: 0, Y: 1},   // N
	{X: 1, Y: 1},   // NE
	{X: 1, Y: 0},   // E
	{X: 1, Y: -1},  // SE
	{X: 0, Y: -1},  // S
	{X: -1, Y: -1}, // SW
	{X: -1, Y: 0},  // W
	{X: -1, Y: 1},  // NW
}

// Vector converts a Direction to a Vector
func (d Bearing) Vector() Vector {
	return bearingVectors[d]
}

// IsCardinal returns if this is a cardinal (NESW)
func (d Bearing) IsCardinal() bool {
	switch d {
	case North:
		fallthrough
	case East:
		fallthrough
	case South:
		fallthrough
	case West:
		return true
	}
	return false
}
