package maths

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 0, Abs(0))
	assert.Equal(t, 1, Abs(1))
	assert.Equal(t, 1, Abs(-1))
}

func TestPmod(t *testing.T) {
	assert.Equal(t, 0, Pmod(0, 0))
	assert.Equal(t, 2, Pmod(6, 4))
	assert.Equal(t, 2, Pmod(-6, 4))
	assert.Equal(t, 4, Pmod(-6, 10))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 500, Max(100, 500))
	assert.Equal(t, 1, Max(-4, 1))
	assert.Equal(t, -2, Max(-4, -2))
}

func TestMin(t *testing.T) {
	assert.Equal(t, 100, Min(100, 500))
	assert.Equal(t, -4, Min(-4, 1))
	assert.Equal(t, -4, Min(-4, -2))
}

func TestRoundUp(t *testing.T) {
	assert.Equal(t, 10, RoundUp(10, 5))
	assert.Equal(t, 12, RoundUp(10, 4))
	assert.Equal(t, -8, RoundUp(-8, 4))
	assert.Equal(t, -4, RoundUp(-7, 4))
}

func TestRoundDown(t *testing.T) {
	assert.Equal(t, 10, RoundDown(10, 5))
	assert.Equal(t, 8, RoundDown(10, 4))
	assert.Equal(t, -8, RoundDown(-8, 4))
	assert.Equal(t, -8, RoundDown(-7, 4))

}
