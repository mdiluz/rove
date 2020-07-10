package rove

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/roveapi"
)

// World describes a self contained universe and everything in it
type World struct {
	// TicksPerDay is the amount of ticks in a single day
	TicksPerDay int `json:"ticks-per-day"`

	// Current number of ticks from the start
	CurrentTicks int `json:"current-ticks"`

	// Rovers is a id->data map of all the rovers in the game
	Rovers map[string]Rover `json:"rovers"`

	// Atlas represends the world map of chunks and tiles
	Atlas atlas.Atlas `json:"atlas"`

	// Commands is the set of currently executing command streams per rover
	CommandQueue map[string]CommandStream `json:"commands"`
	// Incoming represents the set of commands to add to the queue at the end of the current tick
	CommandIncoming map[string]CommandStream `json:"incoming"`

	// Mutex to lock around all world operations
	worldMutex sync.RWMutex
	// Mutex to lock around command operations
	cmdMutex sync.RWMutex
	// Set of possible words to use for names
	words []string
}

var wordsFile = os.Getenv("WORDS_FILE")

// NewWorld creates a new world object
func NewWorld(chunkSize int) *World {

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
		Rovers:          make(map[string]Rover),
		CommandQueue:    make(map[string]CommandStream),
		CommandIncoming: make(map[string]CommandStream),
		Atlas:           atlas.NewChunkAtlas(chunkSize),
		words:           lines,
		TicksPerDay:     24,
		CurrentTicks:    0,
	}
}

// SpawnRover adds an rover to the game
func (w *World) SpawnRover() (string, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	// Initialise the rover
	rover := Rover{
		Range:            4,
		Integrity:        10,
		MaximumIntegrity: 10,
		Capacity:         10,
		Charge:           10,
		MaximumCharge:    10,
		Name:             uuid.New().String(),
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
	rover.Pos = maths.Vector{
		X: 10 - rand.Intn(20),
		Y: 10 - rand.Intn(20),
	}

	// Seach until we error (run out of world)
	for {
		_, obj := w.Atlas.QueryPosition(rover.Pos)
		if !obj.IsBlocking() {
			break
		} else {
			// Try and spawn to the east of the blockage
			rover.Pos.Add(maths.Vector{X: 1, Y: 0})
		}

	}

	// Add a log entry for robot creation
	rover.AddLogEntryf("created at %+v", rover.Pos)

	// Append the rover to the list
	w.Rovers[rover.Name] = rover

	return rover.Name, nil
}

// GetRover gets a specific rover by name
func (w *World) GetRover(rover string) (Rover, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return Rover{}, fmt.Errorf("Failed to find rover with name: %s", rover)
	}
	return i, nil
}

// RoverRecharge charges up a rover
func (w *World) RoverRecharge(rover string) (int, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return 0, fmt.Errorf("Failed to find rover with name: %s", rover)
	}

	// We can only recharge during the day
	if !w.Daytime() {
		return i.Charge, nil
	}

	// Add one charge
	if i.Charge < i.MaximumCharge {
		i.Charge++
		i.AddLogEntryf("recharged to %d", i.Charge)
	}
	w.Rovers[rover] = i

	return i.Charge, nil
}

// RoverBroadcast broadcasts a message to nearby rovers
func (w *World) RoverBroadcast(rover string, message []byte) (err error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return fmt.Errorf("Failed to find rover with name: %s", rover)
	}

	// Use up a charge as needed, if available
	if i.Charge == 0 {
		return
	}
	i.Charge--

	// Check all rovers
	for r, rover := range w.Rovers {
		if rover.Name == i.Name {
			continue
		}

		// Check if this rover is within range
		if i.Pos.Distance(rover.Pos) < float64(i.Range) {
			rover.AddLogEntryf("recieved %s from %s", string(message), i.Name)
			w.Rovers[r] = rover
		}
	}

	i.AddLogEntryf("broadcasted %s", string(message))
	w.Rovers[rover] = i
	return
}

// DestroyRover Removes an rover from the game
func (w *World) DestroyRover(rover string) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	_, ok := w.Rovers[rover]
	if !ok {
		return fmt.Errorf("no rover matching id")
	}

	delete(w.Rovers, rover)
	return nil
}

// RoverPosition returns the position of the rover
func (w *World) RoverPosition(rover string) (maths.Vector, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return maths.Vector{}, fmt.Errorf("no rover matching id")
	}
	return i.Pos, nil
}

// SetRoverPosition sets the position of the rover
func (w *World) SetRoverPosition(rover string, pos maths.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return fmt.Errorf("no rover matching id")
	}

	i.Pos = pos
	w.Rovers[rover] = i
	return nil
}

// RoverInventory returns the inventory of a requested rover
func (w *World) RoverInventory(rover string) ([]atlas.Object, error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return nil, fmt.Errorf("no rover matching id")
	}
	return i.Inventory, nil
}

// WarpRover sets an rovers position
func (w *World) WarpRover(rover string, pos maths.Vector) error {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return fmt.Errorf("no rover matching id")
	}
	// Nothing to do if these positions match
	if i.Pos == pos {
		return nil
	}

	// Check the tile is not blocked
	_, obj := w.Atlas.QueryPosition(pos)
	if obj.IsBlocking() {
		return fmt.Errorf("can't warp rover to occupied tile, check before warping")
	}

	i.Pos = pos
	w.Rovers[rover] = i
	return nil
}

// MoveRover attempts to move a rover in a specific direction
func (w *World) MoveRover(rover string, b maths.Bearing) (maths.Vector, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	i, ok := w.Rovers[rover]
	if !ok {
		return maths.Vector{}, fmt.Errorf("no rover matching id")
	}

	// Ensure the rover has energy
	if i.Charge <= 0 {
		return i.Pos, nil
	}
	i.Charge--

	// Try the new move position
	newPos := i.Pos.Added(b.Vector())

	// Get the tile and verify it's empty
	_, obj := w.Atlas.QueryPosition(newPos)
	if !obj.IsBlocking() {
		i.AddLogEntryf("moved %s to %+v", b.String(), newPos)
		// Perform the move
		i.Pos = newPos
		w.Rovers[rover] = i
	} else {
		// If it is a blocking tile, reduce the rover integrity
		i.AddLogEntryf("tried to move %s to %+v", b.String(), newPos)
		i.Integrity = i.Integrity - 1
		i.AddLogEntryf("had a collision, new integrity %d", i.Integrity)
		if i.Integrity == 0 {
			// TODO: The rover needs to be left dormant with the player
		} else {
			w.Rovers[rover] = i
		}
	}

	return i.Pos, nil
}

// RoverStash will stash an item at the current rovers position
func (w *World) RoverStash(rover string) (atlas.Type, error) {
	w.worldMutex.Lock()
	defer w.worldMutex.Unlock()

	r, ok := w.Rovers[rover]
	if !ok {
		return atlas.ObjectNone, fmt.Errorf("no rover matching id")
	}

	// Can't pick up when full
	if len(r.Inventory) >= r.Capacity {
		return atlas.ObjectNone, nil
	}

	// Ensure the rover has energy
	if r.Charge <= 0 {
		return atlas.ObjectNone, nil
	}
	r.Charge--

	_, obj := w.Atlas.QueryPosition(r.Pos)
	if !obj.IsStashable() {
		return atlas.ObjectNone, nil
	}

	r.AddLogEntryf("stashed %c", obj.Type)
	r.Inventory = append(r.Inventory, obj)
	w.Rovers[rover] = r
	w.Atlas.SetObject(r.Pos, atlas.Object{Type: atlas.ObjectNone})
	return obj.Type, nil
}

// RadarFromRover can be used to query what a rover can currently see
func (w *World) RadarFromRover(rover string) (radar []byte, objs []byte, err error) {
	w.worldMutex.RLock()
	defer w.worldMutex.RUnlock()

	r, ok := w.Rovers[rover]
	if !ok {
		err = fmt.Errorf("no rover matching id")
		return
	}

	// The radar should span in range direction on each axis, plus the row/column the rover is currently on
	radarSpan := (r.Range * 2) + 1
	roverPos := r.Pos

	// Get the radar min and max values
	radarMin := maths.Vector{
		X: roverPos.X - r.Range,
		Y: roverPos.Y - r.Range,
	}
	radarMax := maths.Vector{
		X: roverPos.X + r.Range,
		Y: roverPos.Y + r.Range,
	}

	// Gather up all tiles within the range
	radar = make([]byte, radarSpan*radarSpan)
	objs = make([]byte, radarSpan*radarSpan)
	for j := radarMin.Y; j <= radarMax.Y; j++ {
		for i := radarMin.X; i <= radarMax.X; i++ {
			q := maths.Vector{X: i, Y: j}

			tile, obj := w.Atlas.QueryPosition(q)

			// Get the position relative to the bottom left of the radar
			relative := q.Added(radarMin.Negated())
			index := relative.X + relative.Y*radarSpan
			radar[index] = tile
			objs[index] = byte(obj.Type)
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
			objs[index] = byte(atlas.ObjectRover)
		}
	}

	return radar, objs, nil
}

// RoverCommands returns current commands for the given rover
func (w *World) RoverCommands(rover string) (incoming []Command, queued []Command) {
	if c, ok := w.CommandIncoming[rover]; ok {
		incoming = c
	}
	if c, ok := w.CommandQueue[rover]; ok {
		queued = c
	}
	return
}

// Enqueue will queue the commands given
func (w *World) Enqueue(rover string, commands ...Command) error {

	// First validate the commands
	for _, c := range commands {
		switch c.Command {
		case roveapi.CommandType_move:
			if _, err := maths.FromString(c.Bearing); err != nil {
				return fmt.Errorf("unknown bearing: %s", c.Bearing)
			}
		case roveapi.CommandType_broadcast:
			if len(c.Message) > 3 {
				return fmt.Errorf("too many characters in message (limit 3): %d", len(c.Message))
			}
			for _, b := range c.Message {
				if b < 37 || b > 126 {
					return fmt.Errorf("invalid message character: %c", b)
				}
			}
		case roveapi.CommandType_stash:
		case roveapi.CommandType_repair:
		case roveapi.CommandType_recharge:
			// Nothing to verify
		default:
			return fmt.Errorf("unknown command: %s", c.Command)
		}
	}

	// Lock our commands edit
	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	// Override the incoming command set
	w.CommandIncoming[rover] = commands

	return nil
}

// EnqueueAllIncoming will enqueue the incoming commands
func (w *World) EnqueueAllIncoming() {
	// Add any incoming commands from this tick and clear that queue
	for id, incoming := range w.CommandIncoming {
		commands := w.CommandQueue[id]
		commands = append(commands, incoming...)
		w.CommandQueue[id] = commands
	}
	w.CommandIncoming = make(map[string]CommandStream)
}

// ExecuteCommandQueues will execute any commands in the current command queue
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

	// Increment the current tick count
	w.CurrentTicks++
}

// ExecuteCommand will execute a single command
func (w *World) ExecuteCommand(c *Command, rover string) (err error) {
	log.Printf("Executing command: %+v for %s\n", *c, rover)

	switch c.Command {
	case roveapi.CommandType_move:
		if dir, err := maths.FromString(c.Bearing); err != nil {
			return err
		} else if _, err := w.MoveRover(rover, dir); err != nil {
			return err
		}

	case roveapi.CommandType_stash:
		if _, err := w.RoverStash(rover); err != nil {
			return err
		}

	case roveapi.CommandType_repair:
		r, err := w.GetRover(rover)
		if err != nil {
			return err
		}
		// Consume an inventory item to repair if possible
		if len(r.Inventory) > 0 && r.Integrity < r.MaximumIntegrity {
			r.Inventory = r.Inventory[:len(r.Inventory)-1]
			r.Integrity = r.Integrity + 1
			r.AddLogEntryf("repaired self to %d", r.Integrity)
			w.Rovers[rover] = r
		}
	case roveapi.CommandType_recharge:
		_, err := w.RoverRecharge(rover)
		if err != nil {
			return err
		}
	case roveapi.CommandType_broadcast:
		if err := w.RoverBroadcast(rover, c.Message); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown command: %s", c.Command)
	}

	return
}

// Daytime returns if it's currently daytime
// for simplicity this uses the 1st half of the day as daytime, the 2nd half as nighttime
func (w *World) Daytime() bool {
	tickInDay := w.CurrentTicks % w.TicksPerDay
	return tickInDay < w.TicksPerDay/2
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
