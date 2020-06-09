package vector

import "math"

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