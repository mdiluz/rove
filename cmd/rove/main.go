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
	fmt.Println("\tcommands\tissues commands to the rover")
	fmt.Println("\tradar   \tgathers radar data for the current rover")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

var home = os.Getenv("HOME")
var filepath = path.Join(home, ".local/share/rove.json")

var host = flag.String("host", "api.rove-game.com", "path to game host server")
var dataPath = flag.String("data", filepath, "data file for storage")

// Data is used to store internal data
type Data struct {
	Account string `json:"account,omitempty"`
	Host    string `json:"host,omitempty"`
}

var name = flag.String("name", "", "used with status command for the account name")

func verifyId(d Data) {
	if len(d.Account) == 0 {
		fmt.Fprintf(os.Stderr, "No account ID set, must register first or set \"account\" value in %s\n", *dataPath)
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
	var data = Data{}
	_, err := os.Stat(*dataPath)
	if !os.IsNotExist(err) {
		if b, err := ioutil.ReadFile(*dataPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file %s error: %s\n", *dataPath, err)
			os.Exit(1)
		} else if len(b) == 0 {
			fmt.Fprintf(os.Stderr, "file %s was empty, assumin fresh data\n", *dataPath)

		} else if err := json.Unmarshal(b, &data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal file %s error: %s\n", *dataPath, err)
			os.Exit(1)
		}
	}

	var server = rove.Server(*host)

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
			data.Account = response.Id
		}
	case "spawn":
		verifyId(data)
		d := rove.SpawnData{Id: data.Account}
		if response, err := server.Spawn(d); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		} else if !response.Success {
			fmt.Fprintf(os.Stderr, "Server returned failure: %s\n", response.Error)
			os.Exit(1)

		} else {
			fmt.Printf("Spawned at position %+v\n", response.Position)
		}

	case "commands":
		verifyId(data)
		d := rove.CommandsData{Id: data.Account}

		// TODO: Send real commands in

		if response, err := server.Commands(d); err != nil {
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
		verifyId(data)
		d := rove.RadarData{Id: data.Account}
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

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Save out the persistent file
	if b, err := json.MarshalIndent(data, "", "\t"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal data error: %s\n", err)
		os.Exit(1)
	} else {
		if err := ioutil.WriteFile(*dataPath, b, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save file %s error: %s\n", *dataPath, err)
			os.Exit(1)
		}
	}
}
