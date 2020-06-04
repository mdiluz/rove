package game

// Vector desribes a 3D vector
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Add adds one vector to another
func (v *Vector) Add(v2 Vector) {
	v.X += v2.X
	v.Y += v2.Y
}
