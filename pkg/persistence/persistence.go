package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// dataPath global path for persistence
var dataPath = os.TempDir()

// SetPath sets the persistent path for the data storage
func SetPath(path string) error {
	if info, err := os.Stat(path); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("path for persistence is not directory")
	}
	dataPath = path
	return nil
}

// Converts name to a full path
func jsonPath(name string) string {
	return path.Join(dataPath, fmt.Sprintf("rove-%s.json", name))
}

// Save will serialise the interface into a json file
func Save(name string, data interface{}) error {
	path := jsonPath(name)
	if b, err := json.MarshalIndent(data, "", "  "); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(jsonPath(name), b, os.ModePerm); err != nil {
			return err
		}
	}

	fmt.Printf("Saved %s\n", path)
	return nil
}

// Load will load the interface from the json file
func Load(name string, data interface{}) error {
	path := jsonPath(name)
	// Don't load anything if the file doesn't exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("File %s didn't exist, loading with fresh data\n", path)
		return nil
	}

	// Read and unmarshal the json
	if b, err := ioutil.ReadFile(path); err != nil {
		return err
	} else if len(b) == 0 {
		fmt.Printf("File %s was empty, loading with fresh data\n", path)
		return nil
	} else if err := json.Unmarshal(b, data); err != nil {
		return fmt.Errorf("failed to load file %s error: %s", path, err)
	}

	fmt.Printf("Loaded %s\n", path)
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
				return fmt.Errorf("Incorrect args")
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
