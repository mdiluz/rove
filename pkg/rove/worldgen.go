package rove

import (
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/ojrac/opensimplex-go"
)

// WorldGen describes a world gen algorythm
type WorldGen interface {
	// GetTile generates a tile for a location
	GetTile(v maths.Vector) roveapi.Tile

	// GetObject generates an object for a location
	GetObject(v maths.Vector) Object
}

// NoiseWorldGen returns a noise based world generator
type NoiseWorldGen struct {
	// noise describes the noise function
	noise opensimplex.Noise
}

// NewNoiseWorldGen creates a new noise based world generator
func NewNoiseWorldGen(seed int64) WorldGen {
	return &NoiseWorldGen{
		noise: opensimplex.New(seed),
	}
}

const (
	terrainNoiseScale = 15
	rockNoiseScale    = 3
	partsNoiseScale   = 2
)

// GetTile returns the chosen tile at a location
func (g *NoiseWorldGen) GetTile(v maths.Vector) roveapi.Tile {
	t := g.noise.Eval2(float64(v.X)/terrainNoiseScale, float64(v.Y)/terrainNoiseScale)
	switch {
	case t > 0.5:
		return roveapi.Tile_Gravel
	case t > 0.05:
		return roveapi.Tile_Sand
	default:
		return roveapi.Tile_Rock
	}
}

// GetObject returns the chosen object at a location
func (g *NoiseWorldGen) GetObject(v maths.Vector) (obj Object) {
	r := g.noise.Eval2(float64(v.X)/rockNoiseScale, float64(v.Y)/rockNoiseScale)
	switch {
	// Prioritise rocks
	case r > 0.6:
		obj.Type = roveapi.Object_RockLarge
	case r > 0.5:
		obj.Type = roveapi.Object_RockSmall

	default:
		// Otherwise, try some rover parts
		p := g.noise.Eval2(float64(v.X)/partsNoiseScale, float64(v.Y)/partsNoiseScale)
		switch {
		case p > 0.8:
			obj.Type = roveapi.Object_RoverParts
		}
	}

	return obj
}
