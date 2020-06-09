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
