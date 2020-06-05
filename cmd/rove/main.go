package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mdiluz/rove/pkg/rove"
)

var USAGE = ""

// Command usage
func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]... COMMAND\n", os.Args[0])
	fmt.Println("\nCommands:")
	fmt.Println("\tstatus  \tprints the server status")
	fmt.Println("\tregister\tregisters an account and stores it (use with -name)")
	fmt.Println("\tspawn   \tspawns a rover for the current account")
	fmt.Println("\tcommand \tissues commands to the rover")
	fmt.Println("\tradar   \tgathers radar data for the current rover")
	fmt.Println("\trover   \tgets data for current rover")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

var home = os.Getenv("HOME")
var filepath = path.Join(home, ".local/share/rove.json")

var host = flag.String("host", "", "path to game host server")
var data = flag.String("data", filepath, "data file for storage")

// Config is used to store internal data
type Config struct {
	Account string `json:"account,omitempty"`
	Host    string `json:"host,omitempty"`
}

var name = flag.String("name", "", "used with status command for the account name")

func verifyId(d Config) {
	if len(d.Account) == 0 {
		fmt.Fprintf(os.Stderr, "No account ID set, must register first or set \"account\" value in %s\n", *data)
		os.Exit(1)
	}
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	// Verify we have a single command line arg
	args := flag.Args()
	if len(args) != 1 {
		Usage()
		os.Exit(1)
	}

	// Load in the persistent file
	var config = Config{}
	_, err := os.Stat(*data)
	if !os.IsNotExist(err) {
		if b, err := ioutil.ReadFile(*data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file %s error: %s\n", *data, err)
			os.Exit(1)
		} else if len(b) == 0 {
			fmt.Fprintf(os.Stderr, "file %s was empty, assumin fresh data\n", *data)

		} else if err := json.Unmarshal(b, &config); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal file %s error: %s\n", *data, err)
			os.Exit(1)
		}
	}

	// If there's a host set on the command line, override the one in the config
	if len(*host) != 0 {
		config.Host = *host
	}

	// If there's still no host, bail
	if len(config.Host) == 0 {
		fmt.Fprintln(os.Stderr, "no host set, please set one with -host")
		os.Exit(1)
	}

	// Set up the server
	var server = rove.Server(config.Host)

	// Print the config info
	fmt.Printf("host: %s\taccount: %s\n", config.Host, config.Account)

	// Handle all the commands
	command := args[0]
	switch command {
	case "status":
		if response, err := server.Status(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else {
			fmt.Printf("Ready: %t\n", response.Ready)
			fmt.Printf("Version: %s\n", response.Version)
		}

	case "register":
		d := rove.RegisterData{
			Name: *name,
		}
		if response, err := server.Register(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			fmt.Printf("Registered account with id: %s\n", response.Id)
			config.Account = response.Id
		}
	case "spawn":
		verifyId(config)
		d := rove.SpawnData{Id: config.Account}
		if response, err := server.Spawn(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			fmt.Printf("Spawned at position %+v\n", response.Position)
		}

	case "command":
		verifyId(config)
		d := rove.CommandData{Id: config.Account}

		// TODO: Send real commands in

		if response, err := server.Command(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			// TODO: Pretify the response
			fmt.Printf("%+v\n", response)
		}

	case "radar":
		verifyId(config)
		d := rove.RadarData{Id: config.Account}
		if response, err := server.Radar(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			// TODO: Pretify the response
			fmt.Printf("%+v\n", response)
		}

	case "rover":
		verifyId(config)
		d := rove.RoverData{Id: config.Account}
		if response, err := server.Rover(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			fmt.Printf("position: %v\n", response.Position)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Save out the persistent file
	if b, err := json.MarshalIndent(config, "", "\t"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal data error: %s\n", err)
		os.Exit(1)
	} else {
		if err := ioutil.WriteFile(*data, b, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save file %s error: %s\n", *data, err)
			os.Exit(1)
		}
	}
}
