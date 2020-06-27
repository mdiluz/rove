package game

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/bearing"
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/vector"
)

// World describes a self contained universe and everything in it
type World struct {
	// Rovers is a id->data map of all the rovers in the game
	Rovers map[string]Rover `json:"rovers"`

	// Atlas represends the world map of chunks and tiles
	Atlas atlas.Atlas `json:"atlas"`

	// Mutex to lock around all world operations
	worldMutex sync.RWMutex

	// Commands is the set of currently executing command streams per rover
	CommandQueue map[string]CommandStream `json:"commands"`

	// Incoming represents the set of commands to add to the queue at the end of the current tick
	Incoming map[string]CommandStream `json:"incoming"`

	// Mutex to lock around command operations
	cmdMutex sync.RWMutex

	// Set of possible words to use for names
	words []string
}

var wordsFile = os.Getenv("WORDS_FILE")

// NewWorld creates a new world object
func NewWorld(size, chunkSize int) *World {

	// Try and load the words file
	var lines []string
	if file, err := os.Open(wordsFile); err != nil {
		log.Printf("Couldn't read words file [%s], running without words: %s\n", wordsFile, err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if scanner.Err() != nil {
			log.Printf("Failure during word file scan: %s\n", scanner.Err())
		}
	}

	return &World{
		Rovers:       make(map[string]Rover),
		CommandQueue: make(map[string]CommandStream),
		Incoming:     make(map[string]CommandStream),
		Atlas:        atlas.NewAtlas(size, chunkSize),
		words:        lines,
	}
}

// SpawnWorld spawns a border at the edge of the world atlas
func (w *World) SpawnWorld(fillWorld bool) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()
	if fillWorld {
		if err := w.Atlas.SpawnRocks(); err != nil {
			return err
		}
	}
	return w.Atlas.SpawnWalls()
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover() (string, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	// Initialise the rover
	rover := Rover{
		Range:     4.0,
		Integrity: 10,
		Name:      uuid.New().String(),
	}

	// Assign a random name if we have words
	if len(w.words) > 0 {
		for {
			// Loop until we find a unique name
			name := fmt.Sprintf("%s-%s", w.words[rand.Intn(len(w.words))], w.words[rand.Intn(len(w.words))])
			if _, ok := w.Rovers[name]; !ok {
				rover.Name = name
				break
			}
		}
	}

	// Spawn in a random place near the origin
	rover.Pos = vector.Vector{
		X: w.Atlas.ChunkSize/2 - rand.Intn(w.Atlas.ChunkSize),
		Y: w.Atlas.ChunkSize/2 - rand.Intn(w.Atlas.ChunkSize),
	}

	// Seach until we error (run out of world)
	for {
		if tile, err := w.Atlas.GetTile(rover.Pos); err != nil {
			return "", err
		} else {
			if !objects.IsBlocking(tile) {
				break
			} else {
				// Try and spawn to the east of the blockage
				rover.Pos.Add(vector.Vector{X: 1, Y: 0})
			}
		}
	}

	log.Printf("Spawned rover at %+v\n", rover.Pos)

	// Append the rover to the list
	w.Rovers[rover.Name] = rover

	return rover.Name, nil
}

// GetRover gets a specific rover by name
func (w *World) GetRover(rover string) (Rover, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if i, ok := w.Rovers[rover]; ok {
		return i, nil
	} else {
		return Rover{}, fmt.Errorf("Failed to find rover with name: %s", rover)
	}
}

// Removes an rover from the game
func (w *World) DestroyRover(rover string) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[rover]; ok {
		// Clear the tile
		if err := w.Atlas.SetTile(i.Pos, objects.Empty); err != nil {
			return fmt.Errorf("coudln't clear old rover tile: %s", err)
		}
		delete(w.Rovers, rover)
	} else {
		return fmt.Errorf("no rover matching id")
	}
	return nil
}

// RoverPosition returns the position of the rover
func (w *World) RoverPosition(rover string) (vector.Vector, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if i, ok := w.Rovers[rover]; ok {
		return i.Pos, nil
	} else {
		return vector.Vector{}, fmt.Errorf("no rover matching id")
	}
}

// SetRoverPosition sets the position of the rover
func (w *World) SetRoverPosition(rover string, pos vector.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[rover]; ok {
		i.Pos = pos
		w.Rovers[rover] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// RoverInventory returns the inventory of a requested rover
func (w *World) RoverInventory(rover string) ([]byte, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if i, ok := w.Rovers[rover]; ok {
		return i.Inventory, nil
	} else {
		return nil, fmt.Errorf("no rover matching id")
	}
}

// WarpRover sets an rovers position
func (w *World) WarpRover(rover string, pos vector.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[rover]; ok {
		// Nothing to do if these positions match
		if i.Pos == pos {
			return nil
		}

		// Check the tile is not blocked
		if tile, err := w.Atlas.GetTile(pos); err != nil {
			return fmt.Errorf("coudln't get state of destination rover tile: %s", err)
		} else if objects.IsBlocking(tile) {
			return fmt.Errorf("can't warp rover to occupied tile, check before warping")
		}

		i.Pos = pos
		w.Rovers[rover] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) MoveRover(rover string, b bearing.Bearing) (vector.Vector, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[rover]; ok {
		// Try the new move position
		newPos := i.Pos.Added(b.Vector())

		// Get the tile and verify it's empty
		if tile, err := w.Atlas.GetTile(newPos); err != nil {
			return vector.Vector{}, fmt.Errorf("couldn't get tile for new position: %s", err)
		} else if !objects.IsBlocking(tile) {
			// Perform the move
			i.Pos = newPos
			w.Rovers[rover] = i
		} else {
			// If it is a blocking tile, reduce the rover integrity
			i.Integrity = i.Integrity - 1
			if i.Integrity == 0 {
				// TODO: The rover needs to be left dormant with the player
			} else {
				w.Rovers[rover] = i
			}
		}

		return i.Pos, nil
	} else {
		return vector.Vector{}, fmt.Errorf("no rover matching id")
	}
}

// RoverStash will stash an item at the current rovers position
func (w *World) RoverStash(rover string) (byte, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if r, ok := w.Rovers[rover]; ok {
		if tile, err := w.Atlas.GetTile(r.Pos); err != nil {
			return objects.Empty, err
		} else {
			if objects.IsStashable(tile) {
				r.Inventory = append(r.Inventory, tile)
				w.Rovers[rover] = r
				if err := w.Atlas.SetTile(r.Pos, objects.Empty); err != nil {
					return objects.Empty, err
				} else {
					return tile, nil
				}
			}
		}

	} else {
		return objects.Empty, fmt.Errorf("no rover matching id")
	}

	return objects.Empty, nil
}

// RadarFromRover can be used to query what a rover can currently see
func (w *World) RadarFromRover(rover string) ([]byte, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if r, ok := w.Rovers[rover]; ok {
		// The radar should span in range direction on each axis, plus the row/column the rover is currently on
		radarSpan := (r.Range * 2) + 1
		roverPos := r.Pos

		// Get the radar min and max values
		radarMin := vector.Vector{
			X: roverPos.X - r.Range,
			Y: roverPos.Y - r.Range,
		}
		radarMax := vector.Vector{
			X: roverPos.X + r.Range,
			Y: roverPos.Y + r.Range,
		}

		// Make sure we only query within the actual world
		worldMin, worldMax := w.Atlas.GetWorldExtents()
		scanMin := vector.Vector{
			X: maths.Max(radarMin.X, worldMin.X),
			Y: maths.Max(radarMin.Y, worldMin.Y),
		}
		scanMax := vector.Vector{
			X: maths.Min(radarMax.X, worldMax.X),
			Y: maths.Min(radarMax.Y, worldMax.Y),
		}

		// Gather up all tiles within the range
		var radar = make([]byte, radarSpan*radarSpan)
		for j := scanMin.Y; j <= scanMax.Y; j++ {
			for i := scanMin.X; i <= scanMax.X; i++ {
				q := vector.Vector{X: i, Y: j}

				if tile, err := w.Atlas.GetTile(q); err != nil {
					return nil, fmt.Errorf("failed to query tile: %s", err)

				} else {
					// Get the position relative to the bottom left of the radar
					relative := q.Added(radarMin.Negated())
					index := relative.X + relative.Y*radarSpan
					radar[index] = tile
				}
			}
		}

		// Add all rovers to the radar
		for _, r := range w.Rovers {
			// If the rover is in range
			dist := r.Pos.Added(roverPos.Negated())
			dist = dist.Abs()

			if dist.X <= r.Range && dist.Y <= r.Range {
				relative := r.Pos.Added(radarMin.Negated())
				index := relative.X + relative.Y*radarSpan
				radar[index] = objects.Rover
			}
		}

		// Add this rover
		radar[len(radar)/2] = objects.Rover

		return radar, nil
	} else {
		return nil, fmt.Errorf("no rover matching id")
	}
}

// Enqueue will queue the commands given
func (w *World) Enqueue(rover string, commands ...Command) error {

	// First validate the commands
	for _, c := range commands {
		switch c.Command {
		case CommandMove:
			if _, err := bearing.FromString(c.Bearing); err != nil {
				return fmt.Errorf("unknown bearing: %s", c.Bearing)
			}
		case CommandStash:
		case CommandRepair:
			// Nothing to verify
		default:
			return fmt.Errorf("unknown command: %s", c.Command)
		}
	}

	// Lock our commands edit
	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	// Append the commands to the incoming set
	if cmds, ok := w.Incoming[rover]; ok {
		w.Incoming[rover] = append(cmds, commands...)
	} else {
		w.Incoming[rover] = commands
	}

	return nil
}

// EnqueueAllIncoming will enqueue the incoming commands
func (w *World) EnqueueAllIncoming() {
	// Add any incoming commands from this tick and clear that queue
	for id, incoming := range w.Incoming {
		commands := w.CommandQueue[id]
		commands = append(commands, incoming...)
		w.CommandQueue[id] = commands
	}
	w.Incoming = make(map[string]CommandStream)
}

// Execute will execute any commands in the current command queue
func (w *World) ExecuteCommandQueues() {
	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	// Iterate through all the current commands
	for rover, cmds := range w.CommandQueue {
		if len(cmds) != 0 {
			// Extract the first command in the queue
			c := cmds[0]
			w.CommandQueue[rover] = cmds[1:]

			// Execute the command
			if err := w.ExecuteCommand(&c, rover); err != nil {
				log.Println(err)
				// TODO: Report this error somehow
			}

		} else {
			// Clean out the empty entry
			delete(w.CommandQueue, rover)
		}
	}

	// Add any incoming commands from this tick and clear that queue
	w.EnqueueAllIncoming()
}

// ExecuteCommand will execute a single command
func (w *World) ExecuteCommand(c *Command, rover string) (err error) {
	log.Printf("Executing command: %+v for %s\n", *c, rover)

	switch c.Command {
	case CommandMove:
		if dir, err := bearing.FromString(c.Bearing); err != nil {
			return err
		} else if _, err := w.MoveRover(rover, dir); err != nil {
			return err
		}

	case CommandStash:
		if _, err := w.RoverStash(rover); err != nil {
			return err
		}

	case CommandRepair:
		if r, err := w.GetRover(rover); err != nil {
			return err
		} else {
			// Consume an inventory item to repair
			if len(r.Inventory) > 0 {
				r.Inventory = r.Inventory[:len(r.Inventory)-1]
				r.Integrity = r.Integrity + 1
				w.Rovers[rover] = r
			}
		}
	default:
		return fmt.Errorf("unknown command: %s", c.Command)
	}

	return
}

// PrintTiles simply prints the input tiles directly for debug
func PrintTiles(tiles []byte) {
	num := int(math.Sqrt(float64(len(tiles))))
	for j := num - 1; j >= 0; j-- {
		for i := 0; i < num; i++ {
			fmt.Printf("%c", tiles[i+num*j])
		}
		fmt.Print("\n")
	}
}

// RLock read locks the world
func (w *World) RLock() {
	w.worldMutex.RLock()
	w.cmdMutex.RLock()
}

// RUnlock read unlocks the world
func (w *World) RUnlock() {
	w.worldMutex.RUnlock()
	w.cmdMutex.RUnlock()
}

// Lock locks the world
func (w *World) Lock() {
	w.worldMutex.Lock()
	w.cmdMutex.Lock()
}

// Unlock unlocks the world
func (w *World) Unlock() {
	w.worldMutex.Unlock()
	w.cmdMutex.Unlock()
}
