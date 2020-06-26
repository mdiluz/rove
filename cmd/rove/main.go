package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mdiluz/rove/pkg/bearing"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var home = os.Getenv("HOME")
var defaultDataPath = path.Join(home, ".local/share/")

// Command usage
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [ARGS...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nCommands")
	fmt.Fprintln(os.Stderr, "\tstatus                     prints the server status")
	fmt.Fprintln(os.Stderr, "\tregister NAME              registers an account and stores it (use with -name)")
	fmt.Fprintln(os.Stderr, "\tcommands COMMAND [VAL...]  issues move command to rover")
	fmt.Fprintln(os.Stderr, "\tradar                      gathers radar data for the current rover")
	fmt.Fprintln(os.Stderr, "\trover                      gets data for current rover")
	fmt.Fprintln(os.Stderr, "\tconfig [HOST]              outputs the local config info, optionally sets host")
	fmt.Fprintln(os.Stderr, "\thelp                       outputs this usage information")
	fmt.Fprintln(os.Stderr, "\tversion                    outputs version info")
	fmt.Fprintln(os.Stderr, "\nEnvironment")
	fmt.Fprintln(os.Stderr, "\tROVE_USER_DATA             path to user data, defaults to "+defaultDataPath)
}

const gRPCport = 9090

// Config is used to store internal data
type Config struct {
	Host     string            `json:"host,omitempty"`
	Accounts map[string]string `json:"accounts,omitempty"`
}

// ConfigPath returns the configuration path
func ConfigPath() string {
	// Allow overriding the data path
	var datapath = defaultDataPath
	var override = os.Getenv("ROVE_USER_DATA")
	if len(override) > 0 {
		datapath = override
	}
	datapath = path.Join(datapath, "rove.json")

	return datapath
}

// LoadConfig loads the config from a chosen path
func LoadConfig() (config Config, err error) {

	datapath := ConfigPath()
	config.Accounts = make(map[string]string)

	// Create the path if needed
	path := filepath.Dir(datapath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	} else {
		// Read the file
		_, err = os.Stat(datapath)
		if !os.IsNotExist(err) {
			if b, err := ioutil.ReadFile(datapath); err != nil {
				return Config{}, fmt.Errorf("failed to read file %s error: %s", datapath, err)

			} else if len(b) == 0 {
				return Config{}, fmt.Errorf("file %s was empty, assumin fresh data", datapath)

			} else if err := json.Unmarshal(b, &config); err != nil {
				return Config{}, fmt.Errorf("failed to unmarshal file %s error: %s", datapath, err)

			}
		}
	}

	return
}

// SaveConfig saves the config out
func SaveConfig(config Config) error {
	// Save out the persistent file
	datapath := ConfigPath()
	if b, err := json.MarshalIndent(config, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal data error: %s", err)
	} else if err := ioutil.WriteFile(datapath, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file %s error: %s", datapath, err)
	}

	return nil
}

// verifyID will verify an account ID
func verifyID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("no account ID set, must register first")
	}
	return nil
}

// InnerMain wraps the main function so we can test it
func InnerMain(command string, args ...string) error {

	// Early simple bails
	switch command {
	case "help":
		printUsage()
		return nil
	case "version":
		fmt.Println(version.Version)
		return nil
	}

	// Load in the persistent file
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Run config command before server needed
	if command == "config" {
		if len(args) > 0 {
			config.Host = args[0]
		}
		fmt.Printf("host: %s\taccount: %s\n", config.Host, config.Accounts[config.Host])
		return SaveConfig(config)
	}

	// If there's still no host, bail
	if len(config.Host) == 0 {
		return fmt.Errorf("no host set in %s, set one with '%s config {HOST}'", ConfigPath(), os.Args[0])
	}

	// Set up the server
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Host, gRPCport), grpc.WithInsecure())
	if err != nil {
		return err
	}
	var client = rove.NewRoveClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Handle all the commands
	switch command {
	case "status":
		response, err := client.Status(ctx, &rove.StatusRequest{})
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Ready: %t\n", response.Ready)
			fmt.Printf("Version: %s\n", response.Version)
			fmt.Printf("Tick: %d\n", response.Tick)
			fmt.Printf("Next Tick: %s\n", response.NextTick)
		}

	case "register":
		if len(args) == 0 {
			return fmt.Errorf("must pass name to 'register'")
		}
		name := args[0]
		d := rove.RegisterRequest{
			Name: name,
		}
		_, err := client.Register(ctx, &d)
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Registered account with id: %s\n", name)
			config.Accounts[config.Host] = name
		}

	case "commands":
		if len(args) == 0 {
			return fmt.Errorf("must pass commands to 'commands'")
		}

		// Iterate through each command
		var commands []*rove.Command
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "move":
				i++
				if len(args) == i {
					return fmt.Errorf("move command must be passed bearing")
				} else if _, err := bearing.FromString(args[i]); err != nil {
					return err
				}
				commands = append(commands,
					&rove.Command{
						Command: game.CommandMove,
						Bearing: args[i],
					},
				)
			case "stash":
				commands = append(commands,
					&rove.Command{
						Command: game.CommandStash,
					},
				)
			}
		}

		d := rove.CommandsRequest{
			Account:  config.Accounts[config.Host],
			Commands: commands,
		}

		if err := verifyID(d.Account); err != nil {
			return err
		}

		_, err := client.Commands(ctx, &d)
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Request succeeded\n")
		}

	case "radar":
		dat := rove.RadarRequest{Account: config.Accounts[config.Host]}
		if err := verifyID(dat.Account); err != nil {
			return err
		}

		response, err := client.Radar(ctx, &dat)
		switch {
		case err != nil:
			return err

		default:
			// Print out the radar
			game.PrintTiles(response.Tiles)
		}

	case "rover":
		req := rove.RoverRequest{Account: config.Accounts[config.Host]}
		if err := verifyID(req.Account); err != nil {
			return err
		}
		response, err := client.Rover(ctx, &req)

		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("rover info: %+v\n", response)
		}

	default:
		// Print the usage
		fmt.Fprintf(os.Stderr, "Error: unknown command %s\n", command)
		printUsage()
		os.Exit(1)
	}

	return SaveConfig(config)
}

// Simple main
func main() {
	// Bail without any args
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	// Run the inner main
	if err := InnerMain(os.Args[1], os.Args[2:]...); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
