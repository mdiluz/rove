package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
)

var USAGE = ""

// Command usage
func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]... COMMAND\n", os.Args[0])
	fmt.Println("\nCommands:")
	fmt.Println("\tstatus  \tprints the server status")
	fmt.Println("\tregister\tregisters an account and stores it (use with -name)")
	fmt.Println("\tspawn   \tspawns a rover for the current account")
	fmt.Println("\tmove    \tissues move command to rover")
	fmt.Println("\tradar   \tgathers radar data for the current rover")
	fmt.Println("\trover   \tgets data for current rover")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

var home = os.Getenv("HOME")
var filepath = path.Join(home, ".local/share/rove.json")

// General usage
var ver = flag.Bool("version", false, "Display version number")
var host = flag.String("host", "", "path to game host server")
var data = flag.String("data", filepath, "data file for storage")

// For register command
var name = flag.String("name", "", "used with status command for the account name")

// For the duration command
var duration = flag.Int("duration", 1, "used for the move command duration")
var bearing = flag.String("bearing", "", "used for the move command bearing (compass direction)")

// Config is used to store internal data
type Config struct {
	Host     string            `json:"host,omitempty"`
	Accounts map[string]string `json:"accounts,omitempty"`
}

// verifyId will verify an account ID
func verifyId(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("no account ID set, must register first")
	}
	return nil
}

// InnerMain wraps the main function so we can test it
func InnerMain(command string) error {

	// Load in the persistent file
	var config = Config{
		Accounts: make(map[string]string),
	}
	_, err := os.Stat(*data)
	if !os.IsNotExist(err) {
		if b, err := ioutil.ReadFile(*data); err != nil {
			return fmt.Errorf("failed to read file %s error: %s", *data, err)

		} else if len(b) == 0 {
			return fmt.Errorf("file %s was empty, assumin fresh data", *data)

		} else if err := json.Unmarshal(b, &config); err != nil {
			return fmt.Errorf("failed to unmarshal file %s error: %s", *data, err)

		}
	}

	// If there's a host set on the command line, override the one in the config
	if len(*host) != 0 {
		config.Host = *host
	}

	// If there's still no host, bail
	if len(config.Host) == 0 {
		return fmt.Errorf("no host set, please set one with -host")
	}

	// Set up the server
	var server = rove.Server(config.Host)

	// Grab the account
	var account = config.Accounts[config.Host]

	// Print the config info
	fmt.Printf("host: %s\taccount: %s\n", config.Host, account)

	// Handle all the commands
	switch command {
	case "status":
		if response, err := server.Status(); err != nil {
			return err

		} else {
			fmt.Printf("Ready: %t\n", response.Ready)
			fmt.Printf("Version: %s\n", response.Version)
			fmt.Printf("Tick: %d\n", response.Tick)
			fmt.Printf("Next Tick: %s\n", response.NextTick)
		}

	case "register":
		d := rove.RegisterData{
			Name: *name,
		}
		if response, err := server.Register(d); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			fmt.Printf("Registered account with id: %s\n", response.Id)
			config.Accounts[config.Host] = response.Id
		}
	case "spawn":
		d := rove.SpawnData{}
		if err := verifyId(account); err != nil {
			return err
		} else if response, err := server.Spawn(account, d); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			fmt.Printf("Spawned rover with attributes %+v\n", response.Attributes)
		}

	case "move":
		d := rove.CommandData{
			Commands: []game.Command{
				{
					Command:  game.CommandMove,
					Duration: *duration,
					Bearing:  *bearing,
				},
			},
		}

		if err := verifyId(account); err != nil {
			return err
		} else if response, err := server.Command(account, d); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			// TODO: Pretify the response
			fmt.Printf("%+v\n", response)
		}

	case "radar":
		if err := verifyId(account); err != nil {
			return err
		} else if response, err := server.Radar(account); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			fmt.Printf("nearby rovers: %+v\n", response.Rovers)
		}

	case "rover":
		if err := verifyId(account); err != nil {
			return err
		} else if response, err := server.Rover(account); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			fmt.Printf("attributes: %+v\n", response.Attributes)
		}

	default:
		return fmt.Errorf("Unknown command: %s", command)
	}

	// Save out the persistent file
	if b, err := json.MarshalIndent(config, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal data error: %s", err)
	} else {
		if err := ioutil.WriteFile(*data, b, os.ModePerm); err != nil {
			return fmt.Errorf("failed to save file %s error: %s", *data, err)
		}
	}

	return nil
}

// Simple main
func main() {
	flag.Usage = Usage
	flag.Parse()

	// Print the version if requested
	if *ver {
		fmt.Println(version.Version)
		return
	}

	// Verify we have a single command line arg
	args := flag.Args()
	if len(args) != 1 {
		Usage()
		os.Exit(1)
	}

	// Run the inner main
	if err := InnerMain(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
