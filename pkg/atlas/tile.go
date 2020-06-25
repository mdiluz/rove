package atlas

const (
	TileEmpty     = byte(' ')
	TileRover     = byte('R')
	TileSmallRock = byte('o')
	TileLargeRock = byte('O')
)

// BlockingTiles describes any tiles that block
var BlockingTiles = [...]byte{
	TileLargeRock,
}

// Check if a tile is a blocking tile
func IsBlocking(tile byte) bool {
	for _, t := range BlockingTiles {
		if tile == t {
			return true
		}
	}
	return false
}
