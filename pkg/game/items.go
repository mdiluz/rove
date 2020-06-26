package game

// Each item is a specific type
const (
	ItemNone = byte(0)

	// Describes a single rock
	ItemRock = byte(1)
)

// Item describes an item that can be held
type Item struct {
	Type byte `json:"type"`
}
