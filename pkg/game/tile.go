package game

// Tile represents the type of a tile on the map
type Tile byte

const (
	TileEmpty = Tile(0)
	TileRover = Tile(1)

	// TODO: Is there even a difference between these two?
	TileWall = Tile(2)
	TileRock = Tile(3)
)
