package game

// Tile represents a single tile on the map
type Tile struct {
	// Kind represends the kind of tile this is
	Kind int `json:"kind"`
}

const (
	ChunkDimensions = 10
)

// Chunk represents a fixed square grid of tiles
type Chunk struct {
	// Tiles represents the tiles within the chunk
	Tiles [ChunkDimensions][ChunkDimensions]Tile `json:"tiles"`
}

const (
	// Use a fixed map dimension for now
	AtlasDimensions = 10
)

// Atlas represents a grid of Chunks
// TODO: Make this resizable
type Atlas struct {
	Chunks [AtlasDimensions][AtlasDimensions]Chunk `json:"chunks"`
}
