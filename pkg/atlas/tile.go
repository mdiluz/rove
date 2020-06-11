package atlas

// Tile represents the type of a tile on the map
type Tile byte

const (
	TileEmpty = Tile(0)
	TileRover = Tile(1)

	TileWall = Tile(2)
	TileRock = Tile(3)
)
