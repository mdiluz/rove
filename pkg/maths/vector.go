package maths

import (
	"math"

	"github.com/mdiluz/rove/proto/roveapi"
)

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

// DividedFloor returns the vector divided but floors the value regardless
func (v Vector) DividedFloor(val int) Vector {
	x := float64(v.X) / float64(val)

	if x < 0 {
		x = math.Floor(x)
	} else {
		x = math.Floor(x)
	}
	y := float64(v.Y) / float64(val)
	if y < 0 {
		y = math.Floor(y)
	} else {
		y = math.Floor(y)
	}

	return Vector{X: int(x), Y: int(y)}
}

// Abs returns an absolute version of the vector
func (v Vector) Abs() Vector {
	return Vector{Abs(v.X), Abs(v.Y)}
}

// Min2 returns the minimum values in both vectors
func Min2(v1 Vector, v2 Vector) Vector {
	return Vector{Min(v1.X, v2.X), Min(v1.Y, v2.Y)}
}

// Max2 returns the max values in both vectors
func Max2(v1 Vector, v2 Vector) Vector {
	return Vector{Max(v1.X, v2.X), Max(v1.Y, v2.Y)}
}

// BearingToVector converts a bearing to a vector
func BearingToVector(b roveapi.Bearing) Vector {
	switch b {
	case roveapi.Bearing_North:
		return Vector{Y: 1}
	case roveapi.Bearing_East:
		return Vector{X: 1}
	case roveapi.Bearing_South:
		return Vector{Y: -1}
	case roveapi.Bearing_West:
		return Vector{X: -1}
	}

	return Vector{}
}
