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
	"github.com/mdiluz/rove/pkg/vector"
)

// World describes a self contained universe and everything in it
type World struct {
	// Rovers is a id->data map of all the rovers in the game
	Rovers map[uuid.UUID]Rover `json:"rovers"`

	// Atlas represends the world map of chunks and tiles
	Atlas atlas.Atlas `json:"atlas"`

	// Mutex to lock around all world operations
	worldMutex sync.RWMutex

	// Commands is the set of currently executing command streams per rover
	CommandQueue map[uuid.UUID]CommandStream `json:"commands"`

	// Incoming represents the set of commands to add to the queue at the end of the current tick
	Incoming map[uuid.UUID]CommandStream `json:"incoming"`

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
		Rovers:       make(map[uuid.UUID]Rover),
		CommandQueue: make(map[uuid.UUID]CommandStream),
		Incoming:     make(map[uuid.UUID]CommandStream),
		Atlas:        atlas.NewAtlas(size, chunkSize),
		words:        lines,
	}
}

// SpawnWorld spawns a border at the edge of the world atlas
func (w *World) SpawnWorld(fillWorld bool) error {
	if fillWorld {
		if err := w.Atlas.SpawnRocks(); err != nil {
			return err
		}
	}
	return w.Atlas.SpawnWalls()
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover() (uuid.UUID, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	// Initialise the rover
	rover := Rover{
		Id: uuid.New(),
		Attributes: RoverAttributes{
			Speed:    1.0,
			Range:    5.0,
			Capacity: 5,
			Name:     "rover",
		},
	}

	// Assign a random name if we have words
	if len(w.words) > 0 {
		rover.Attributes.Name = fmt.Sprintf("%s-%s", w.words[rand.Intn(len(w.words))], w.words[rand.Intn(len(w.words))])
	}

	// Spawn in a random place near the origin
	rover.Pos = vector.Vector{
		X: w.Atlas.ChunkSize/2 - rand.Intn(w.Atlas.ChunkSize),
		Y: w.Atlas.ChunkSize/2 - rand.Intn(w.Atlas.ChunkSize),
	}

	// Seach until we error (run out of world)
	for {
		if tile, err := w.Atlas.GetTile(rover.Pos); err != nil {
			return uuid.Nil, err
		} else {
			if !atlas.IsBlocking(tile) {
				break
			} else {
				// Try and spawn to the east of the blockage
				rover.Pos.Add(vector.Vector{X: 1, Y: 0})
			}
		}
	}

	log.Printf("Spawned rover at %+v\n", rover.Pos)

	// Append the rover to the list
	w.Rovers[rover.Id] = rover

	return rover.Id, nil
}

// Removes an rover from the game
func (w *World) DestroyRover(id uuid.UUID) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		// Clear the tile
		if err := w.Atlas.SetTile(i.Pos, atlas.TileEmpty); err != nil {
			return fmt.Errorf("coudln't clear old rover tile: %s", err)
		}
		delete(w.Rovers, id)
	} else {
		return fmt.Errorf("no rover matching id")
	}
	return nil
}

// RoverPosition returns the position of the rover
func (w *World) RoverPosition(id uuid.UUID) (vector.Vector, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if i, ok := w.Rovers[id]; ok {
		return i.Pos, nil
	} else {
		return vector.Vector{}, fmt.Errorf("no rover matching id")
	}
}

// SetRoverPosition sets the position of the rover
func (w *World) SetRoverPosition(id uuid.UUID, pos vector.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		i.Pos = pos
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// RoverAttributes returns the attributes of a requested rover
func (w *World) RoverAttributes(id uuid.UUID) (RoverAttributes, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if i, ok := w.Rovers[id]; ok {
		return i.Attributes, nil
	} else {
		return RoverAttributes{}, fmt.Errorf("no rover matching id")
	}
}

// SetRoverAttributes sets the attributes of a requested rover
func (w *World) SetRoverAttributes(id uuid.UUID, attributes RoverAttributes) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		i.Attributes = attributes
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// WarpRover sets an rovers position
func (w *World) WarpRover(id uuid.UUID, pos vector.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		// Nothing to do if these positions match
		if i.Pos == pos {
			return nil
		}

		// Check the tile is not blocked
		if tile, err := w.Atlas.GetTile(pos); err != nil {
			return fmt.Errorf("coudln't get state of destination rover tile: %s", err)
		} else if atlas.IsBlocking(tile) {
			return fmt.Errorf("can't warp rover to occupied tile, check before warping")
		}

		i.Pos = pos
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) MoveRover(id uuid.UUID, b bearing.Bearing) (vector.Vector, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		// Calculate the distance
		distance := i.Attributes.Speed

		// Calculate the full movement based on the bearing
		move := b.Vector().Multiplied(distance)

		// Try the new move position
		newPos := i.Pos.Added(move)

		// Get the tile and verify it's empty
		if tile, err := w.Atlas.GetTile(newPos); err != nil {
			return vector.Vector{}, fmt.Errorf("couldn't get tile for new position: %s", err)
		} else if !atlas.IsBlocking(tile) {
			// Perform the move
			i.Pos = newPos
			w.Rovers[id] = i
		}

		return i.Pos, nil
	} else {
		return vector.Vector{}, fmt.Errorf("no rover matching id")
	}
}

// RadarFromRover can be used to query what a rover can currently see
func (w *World) RadarFromRover(id uuid.UUID) ([]byte, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if r, ok := w.Rovers[id]; ok {
		// The radar should span in range direction on each axis, plus the row/column the rover is currently on
		radarSpan := (r.Attributes.Range * 2) + 1
		roverPos := r.Pos

		// Get the radar min and max values
		radarMin := vector.Vector{
			X: roverPos.X - r.Attributes.Range,
			Y: roverPos.Y - r.Attributes.Range,
		}
		radarMax := vector.Vector{
			X: roverPos.X + r.Attributes.Range,
			Y: roverPos.Y + r.Attributes.Range,
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

			if dist.X <= r.Attributes.Range && dist.Y <= r.Attributes.Range {
				relative := r.Pos.Added(radarMin.Negated())
				index := relative.X + relative.Y*radarSpan
				radar[index] = atlas.TileRover
			}
		}

		// Add this rover
		radar[len(radar)/2] = atlas.TileRover

		return radar, nil
	} else {
		return nil, fmt.Errorf("no rover matching id")
	}
}

// Enqueue will queue the commands given
func (w *World) Enqueue(rover uuid.UUID, commands ...Command) error {

	// First validate the commands
	for _, c := range commands {
		switch c.Command {
		case "move":
			if _, err := bearing.FromString(c.Bearing); err != nil {
				return fmt.Errorf("unknown bearing: %s", c.Bearing)
			}
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
	w.Incoming = make(map[uuid.UUID]CommandStream)
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

			// Execute the command and clear up if requested
			if done, err := w.ExecuteCommand(&c, rover); err != nil {
				w.CommandQueue[rover] = cmds[1:]
				log.Println(err)
			} else if done {
				w.CommandQueue[rover] = cmds[1:]
			} else {
				w.CommandQueue[rover][0] = c
			}

			// If there was an error

		} else {
			// Clean out the empty entry
			delete(w.CommandQueue, rover)
		}
	}

	// Add any incoming commands from this tick and clear that queue
	w.EnqueueAllIncoming()
}

// ExecuteCommand will execute a single command
func (w *World) ExecuteCommand(c *Command, rover uuid.UUID) (finished bool, err error) {
	log.Printf("Executing command: %+v\n", *c)

	switch c.Command {
	case "move":
		if dir, err := bearing.FromString(c.Bearing); err != nil {
			return true, fmt.Errorf("unknown bearing in command %+v, skipping: %s\n", c, err)

		} else if _, err := w.MoveRover(rover, dir); err != nil {
			return true, fmt.Errorf("error moving rover in command %+v, skipping: %s\n", c, err)

		} else {
			// If we've successfully moved, reduce the duration by 1
			c.Duration -= 1

			// If we've used up the full duration, remove it, otherwise update
			if c.Duration == 0 {
				finished = true
			}
		}
	default:
		return true, fmt.Errorf("unknown command: %s", c.Command)
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
