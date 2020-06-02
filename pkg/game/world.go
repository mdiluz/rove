package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
)

// World describes a self contained universe and everything in it
type World struct {
	Instances map[uuid.UUID]Instance `json:"instances"`

	// dataPath is the location for the data to be stored
	dataPath string
}

// Instance describes a single entity or instance of an entity in the world
type Instance struct {
	// id is a unique ID for this instance
	id uuid.UUID

	// pos represents where this instance is in the world
	pos Position
}

const kWorldFileName = "rove-world.json"

// NewWorld creates a new world object
func NewWorld(data string) *World {
	return &World{
		Instances: make(map[uuid.UUID]Instance),
		dataPath:  data,
	}
}

// path returns the full path to the data file
func (w *World) path() string {
	return path.Join(w.dataPath, kWorldFileName)
}

// Load will load the accountant from data
func (w *World) Load() error {
	// Don't load anything if the file doesn't exist
	_, err := os.Stat(w.path())
	if os.IsNotExist(err) {
		fmt.Printf("File %s didn't exist, loading with fresh world data\n", w.path())
		return nil
	}

	if b, err := ioutil.ReadFile(w.path()); err != nil {
		return err
	} else if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	return nil
}

// Save will save the accountant data out
func (w *World) Save() error {
	if b, err := json.MarshalIndent(w, "", "\t"); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(w.path(), b, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// Adds an instance to the game
func (w *World) CreateInstance() uuid.UUID {
	id := uuid.New()

	// Initialise the instance
	instance := Instance{
		id: id,
	}

	// Append the instance to the list
	w.Instances[id] = instance

	return id
}

// GetPosition returns the position of a given instance
func (w World) GetPosition(id uuid.UUID) (Position, error) {
	if i, ok := w.Instances[id]; ok {
		return i.pos, nil
	} else {
		return Position{}, fmt.Errorf("no instance matching id")
	}
}
