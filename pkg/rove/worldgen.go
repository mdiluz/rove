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
	terrainNoiseScale = 6
	rockNoiseScale    = 3
	dormantRoverScale = 25
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
	o := g.noise.Eval2(float64(v.X)/rockNoiseScale, float64(v.Y)/rockNoiseScale)
	switch {
	case o > 0.6:
		obj.Type = roveapi.Object_RockLarge
	case o > 0.5:
		obj.Type = roveapi.Object_RockSmall
	}

	// Very rarely spawn a dormant rover
	if obj.Type == roveapi.Object_ObjectUnknown {
		o = g.noise.Eval2(float64(v.X)/dormantRoverScale, float64(v.Y)/dormantRoverScale)
		if o > 0.8 {
			obj.Type = roveapi.Object_RoverDormant
		}
	}

	return obj
}
