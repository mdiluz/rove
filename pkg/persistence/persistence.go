package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// dataPath global path for persistence
var dataPath = os.TempDir()

// SetPath sets the persistent path for the data storage
func SetPath(p string) error {
	if info, err := os.Stat(p); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("path for persistence is not directory")
	}
	dataPath = p
	return nil
}

// Converts name to a full path
func jsonPath(name string) string {
	return path.Join(dataPath, fmt.Sprintf("rove-%s.json", name))
}

// Save will serialise the interface into a json file
func Save(name string, data interface{}) error {
	p := jsonPath(name)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(p, b, os.ModePerm); err != nil {
		return err
	}

	log.Printf("Saved %s\n", p)
	return nil
}

// Load will load the interface from the json file
func Load(name string, data interface{}) error {
	p := jsonPath(name)
	// Don't load anything if the file doesn't exist
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		log.Printf("File %s didn't exist, loading with fresh data\n", p)
		return nil
	}

	// Read and unmarshal the json
	if b, err := ioutil.ReadFile(p); err != nil {
		return err
	} else if len(b) == 0 {
		log.Printf("File %s was empty, loading with fresh data\n", p)
		return nil
	} else if err := json.Unmarshal(b, data); err != nil {
		return fmt.Errorf("failed to load file %s error: %s", p, err)
	}

	log.Printf("Loaded %s\n", p)
	return nil
}

// saveLoadFunc defines a type of function to save or load an interface
type saveLoadFunc func(string, interface{}) error

func doAll(f saveLoadFunc, args ...interface{}) error {
	var name string
	for i, a := range args {
		if i%2 == 0 {
			var ok bool
			name, ok = a.(string)
			if !ok {
				return fmt.Errorf("incorrect args")
			}
		} else {
			if err := f(name, a); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveAll allows for saving multiple structures in a single call
func SaveAll(args ...interface{}) error {
	return doAll(Save, args...)
}

// LoadAll allows for loading multiple structures in a single call
func LoadAll(args ...interface{}) error {
	return doAll(Load, args...)
}
