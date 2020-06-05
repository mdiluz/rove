package game

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type VectorTest struct {
	name string
	vec  Vector
	arg  Vector
	want Vector
}

var additionTests = []VectorTest{
	{
		name: "Basic addition 1",
		vec: Vector{
			X: 0.0,
			Y: 0.0,
		},
		arg: Vector{
			1.0,
			1.0,
		},
		want: Vector{
			1.0,
			1.0,
		},
	},
	{
		name: "Basic addition 2",
		vec: Vector{
			X: 1.0,
			Y: 2.0,
		},
		arg: Vector{
			3.0,
			4.0,
		},
		want: Vector{
			4.0,
			6.0,
		},
	},
}

func TestVector_Add(t *testing.T) {

	for _, tt := range additionTests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.vec.X,
				Y: tt.vec.Y,
			}
			v.Add(tt.arg)
			assert.Equal(t, tt.want, v, "Add did not produce expected result")
		})
	}
}

func TestVector_Added(t *testing.T) {
	for _, tt := range additionTests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.vec.X,
				Y: tt.vec.Y,
			}
			assert.Equal(t, tt.want, v.Added(tt.arg), "Added didn't return expected value")
		})
	}
}

func TestVector_Negated(t *testing.T) {
	tests := []struct {
		name string
		vec  Vector
		want Vector
	}{
		{
			name: "Simple check 1",
			vec:  Vector{1, 1},
			want: Vector{-1, -1},
		},
		{
			name: "Simple check 2",
			vec:  Vector{1, -1},
			want: Vector{-1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.vec.X,
				Y: tt.vec.Y,
			}
			if got := v.Negated(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector.Negated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Length(t *testing.T) {
	tests := []struct {
		name string
		vec  Vector
		want float64
	}{
		{
			name: "Simple length 1",
			vec:  Vector{1, 0},
			want: 1,
		}, {
			name: "Simple length 2",
			vec:  Vector{1, -1},
			want: math.Sqrt(2),
		}, {
			name: "Simple length 3",
			vec:  Vector{0, 0},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.vec.X,
				Y: tt.vec.Y,
			}
			if got := v.Length(); got != tt.want {
				t.Errorf("Vector.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Distance(t *testing.T) {
	tests := []struct {
		name string
		vec  Vector
		arg  Vector
		want float64
	}{
		{
			name: "Simple distance 1",
			vec:  Vector{0, 0},
			arg:  Vector{1, 0},
			want: 1,
		},
		{
			name: "Simple distance 2",
			vec:  Vector{1, 1},
			arg:  Vector{-1, -1},
			want: math.Sqrt(8),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.vec.X,
				Y: tt.vec.Y,
			}
			if got := v.Distance(tt.arg); got != tt.want {
				t.Errorf("Vector.Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
