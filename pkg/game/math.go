package game

import (
	"fmt"
	"math"
	"strings"
)

// Abs gets the absolute value of an int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// pmod is a mositive modulo
// golang's % is a "remainder" function si misbehaves for negative modulus inputs
func Pmod(x, d int) int {
	x = x % d
	if x >= 0 {
		return x
	} else if d < 0 {
		return x - d
	} else {
		return x + d
	}
}

// Max returns the highest int
func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the lowest int
func Min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

// Vector desribes a 3D vector
type Vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Add adds one vector to another
func (v *Vector) Add(v2 Vector) {
	v.X += v2.X
	v.Y += v2.Y
}

// Added calculates a new vector
func (v Vector) Added(v2 Vector) Vector {
	v.Add(v2)
	return v
}

// Negated returns a negated vector
func (v Vector) Negated() Vector {
	return Vector{-v.X, -v.Y}
}

// Length returns the length of the vector
func (v Vector) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

// Distance returns the distance between two vectors
func (v Vector) Distance(v2 Vector) float64 {
	// Negate the two vectors and calciate the length
	return v.Added(v2.Negated()).Length()
}

// Multiplied returns the vector multiplied by an int
func (v Vector) Multiplied(val int) Vector {
	return Vector{v.X * val, v.Y * val}
}

// Divided returns the vector divided by an int
func (v Vector) Divided(val int) Vector {
	return Vector{v.X / val, v.Y / val}
}

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

var DirectionVectors = []Vector{
	{0, 1},  // N
	{1, 1},  // NE
	{1, 0},  // E
	{1, -1}, // SE
	{0, -1}, // S
	{-1, 1}, // SW
	{-1, 0}, // W
	{-1, 1}, // NW
}

// Vector converts a Direction to a Vector
func (d Direction) Vector() Vector {
	return DirectionVectors[d]
}
