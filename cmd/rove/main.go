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
	fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [OPTIONS]...\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "\tstatus  \tprints the server status")
	fmt.Fprintln(os.Stderr, "\tregister\tregisters an account and stores it (use with -name)")
	fmt.Fprintln(os.Stderr, "\tspawn   \tspawns a rover for the current account")
	fmt.Fprintln(os.Stderr, "\tmove    \tissues move command to rover")
	fmt.Fprintln(os.Stderr, "\tradar   \tgathers radar data for the current rover")
	fmt.Fprintln(os.Stderr, "\trover   \tgets data for current rover")
	fmt.Fprintln(os.Stderr, "\nOptions:")
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
			fmt.Printf("Request succeeded\n")
		}

	case "radar":
		if err := verifyId(account); err != nil {
			return err
		} else if response, err := server.Radar(account); err != nil {
			return err

		} else if !response.Success {
			return fmt.Errorf("Server returned failure: %s", response.Error)

		} else {
			// Print out the radar
			game.PrintTiles(response.Tiles)
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
		// Print the usage
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		Usage()
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
	flag.CommandLine.Parse(os.Args[2:])

	// Print the version if requested
	if *ver {
		fmt.Println(version.Version)
		return
	}

	// Run the inner main
	if err := InnerMain(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
