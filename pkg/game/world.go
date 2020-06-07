package game

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/tjarratt/babble"
)

// World describes a self contained universe and everything in it
type World struct {
	// Rovers is a id->data map of all the rovers in the game
	Rovers map[uuid.UUID]Rover `json:"rovers"`

	// Atlas represends the world map of chunks and tiles
	Atlas Atlas `json:"atlas"`

	// Mutex to lock around all world operations
	worldMutex sync.RWMutex

	// Commands is the set of currently executing command streams per rover
	CommandQueue map[uuid.UUID]CommandStream `json:"commands"`

	// Mutex to lock around command operations
	cmdMutex sync.RWMutex
}

// NewWorld creates a new world object
func NewWorld() *World {
	return &World{
		Rovers:       make(map[uuid.UUID]Rover),
		CommandQueue: make(map[uuid.UUID]CommandStream),
		Atlas:        NewAtlas(4, 8), // TODO: Choose an appropriate world size
	}
}

// SpawnWorld spawns a border at the edge of the world atlas
func (w *World) SpawnWorld() error {
	return w.Atlas.SpawnWorld()
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover() (uuid.UUID, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	// Initialise the rover
	rover := Rover{
		Id: uuid.New(),
		Attributes: RoverAttributes{

			Speed: 1.0,
			Range: 5.0,

			// Set the name randomly
			Name: babble.NewBabbler().Babble(),
		},
	}

	// Spawn in a random place near the origin
	rover.Attributes.Pos = Vector{
		w.Atlas.ChunkSize - (rand.Int() % (w.Atlas.ChunkSize * 2)),
		w.Atlas.ChunkSize - (rand.Int() % (w.Atlas.ChunkSize * 2)),
	}

	// Seach until we error (run out of world)
	for {
		if tile, err := w.Atlas.GetTile(rover.Attributes.Pos); err != nil {
			return uuid.Nil, err
		} else {
			if tile == TileEmpty {
				break
			} else {
				// Try and spawn to the east of the blockage
				rover.Attributes.Pos.Add(Vector{1, 0})
			}
		}
	}

	// Set the world tile to a rover
	if err := w.Atlas.SetTile(rover.Attributes.Pos, TileRover); err != nil {
		return uuid.Nil, err
	}

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
		if err := w.Atlas.SetTile(i.Attributes.Pos, TileEmpty); err != nil {
			return fmt.Errorf("coudln't clear old rover tile: %s", err)
		}
		delete(w.Rovers, id)
	} else {
		return fmt.Errorf("no rover matching id")
	}
	return nil
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

// WarpRover sets an rovers position
func (w *World) WarpRover(id uuid.UUID, pos Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		// Update the world tile
		// TODO: Make this (and other things) transactional
		if err := w.Atlas.SetTile(pos, TileRover); err != nil {
			return fmt.Errorf("coudln't set rover tile: %s", err)
		} else if err := w.Atlas.SetTile(i.Attributes.Pos, TileEmpty); err != nil {
			return fmt.Errorf("coudln't clear old rover tile: %s", err)
		}

		i.Attributes.Pos = pos
		w.Rovers[id] = i
		return nil
	} else {
		return fmt.Errorf("no rover matching id")
	}
}

// SetPosition sets an rovers position
func (w *World) MoveRover(id uuid.UUID, bearing Direction) (RoverAttributes, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	if i, ok := w.Rovers[id]; ok {
		// Calculate the distance
		distance := i.Attributes.Speed

		// Calculate the full movement based on the bearing
		move := bearing.Vector().Multiplied(distance)

		// Try the new move position
		newPos := i.Attributes.Pos.Added(move)

		// Get the tile and verify it's empty
		if tile, err := w.Atlas.GetTile(newPos); err != nil {
			return i.Attributes, fmt.Errorf("couldn't get tile for new position: %s", err)
		} else if tile == TileEmpty {
			// Set the world tiles
			// TODO: Make this (and other things) transactional
			if err := w.Atlas.SetTile(newPos, TileRover); err != nil {
				return i.Attributes, fmt.Errorf("coudln't set rover tile: %s", err)
			} else if err := w.Atlas.SetTile(i.Attributes.Pos, TileEmpty); err != nil {
				return i.Attributes, fmt.Errorf("coudln't clear old rover tile: %s", err)
			}

			// Perform the move
			i.Attributes.Pos = newPos
			w.Rovers[id] = i
		}

		return i.Attributes, nil
	} else {
		return RoverAttributes{}, fmt.Errorf("no rover matching id")
	}
}

// RadarBlip represents a single blip on the radar
type RadarBlip struct {
	Position Vector `json:"position"`
	Tile     Tile   `json:"tile"`
}

// RadarFromRover can be used to query what a rover can currently see
func (w *World) RadarFromRover(id uuid.UUID) ([]RadarBlip, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	if r, ok := w.Rovers[id]; ok {
		var blips []RadarBlip

		extent := w.Atlas.GetWorldExtent()

		// Get min and max extents to query
		min := Vector{
			Max(-extent, r.Attributes.Pos.X-r.Attributes.Range),
			Max(-extent, r.Attributes.Pos.Y-r.Attributes.Range),
		}
		max := Vector{
			Min(extent-1, r.Attributes.Pos.X+r.Attributes.Range),
			Min(extent-1, r.Attributes.Pos.Y+r.Attributes.Range),
		}

		// Gather up all tiles within the range
		for i := min.X; i < max.X; i++ {
			for j := min.Y; j < max.Y; j++ {

				// Skip this rover
				q := Vector{i, j}
				if q == r.Attributes.Pos {
					continue
				}

				if tile, err := w.Atlas.GetTile(q); err != nil {
					return blips, fmt.Errorf("failed to query tile: %s", err)

				} else if tile != TileEmpty {
					blips = append(blips, RadarBlip{Position: q, Tile: tile})
				}
			}
		}

		return blips, nil
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
			if _, err := DirectionFromString(c.Bearing); err != nil {
				return fmt.Errorf("unknown direction: %s", c.Bearing)
			}
		default:
			return fmt.Errorf("unknown command: %s", c.Command)
		}
	}

	// Lock our commands edit
	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	// Append the commands to the current set
	cmds := w.CommandQueue[rover]
	w.CommandQueue[rover] = append(cmds, commands...)

	return nil
}

// Execute will execute any commands in the current command queue
func (w *World) ExecuteCommandQueues() {
	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	// Iterate through all commands
	for rover, cmds := range w.CommandQueue {
		if len(cmds) != 0 {
			// Extract the first command in the queue
			c := cmds[0]

			// Execute the command and clear up if requested
			if done, err := w.ExecuteCommand(&c, rover); err != nil {
				w.CommandQueue[rover] = cmds[1:]
				fmt.Println(err)
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
}

// ExecuteCommand will execute a single command
func (w *World) ExecuteCommand(c *Command, rover uuid.UUID) (finished bool, err error) {
	fmt.Printf("Executing command: %+v\n", *c)

	switch c.Command {
	case "move":
		if dir, err := DirectionFromString(c.Bearing); err != nil {
			return true, fmt.Errorf("unknown direction in command %+v, skipping: %s\n", c, err)

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
