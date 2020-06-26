package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Command usage
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [OPTIONS]...\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "\tstatus  \tprints the server status")
	fmt.Fprintln(os.Stderr, "\tregister\tregisters an account and stores it (use with -name)")
	fmt.Fprintln(os.Stderr, "\tmove    \tissues move command to rover")
	fmt.Fprintln(os.Stderr, "\tradar   \tgathers radar data for the current rover")
	fmt.Fprintln(os.Stderr, "\trover   \tgets data for current rover")
	fmt.Fprintln(os.Stderr, "\tconfig  \toutputs the local config info")
	fmt.Fprintln(os.Stderr, "\tversion  \toutputs version info")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}

var home = os.Getenv("HOME")
var defaultDataPath = path.Join(home, ".local/share/")

const gRPCport = 9090

// General usage
var host = flag.String("host", "", "path to game host server")
var data = flag.String("data", defaultDataPath, "data location for storage (or $USER_DATA if set)")

// For register command
var name = flag.String("name", "", "used with status command for the account name")

// For the move command
var bearing = flag.String("bearing", "", "used for the move command bearing (compass direction)")

// Config is used to store internal data
type Config struct {
	Host     string            `json:"host,omitempty"`
	Accounts map[string]string `json:"accounts,omitempty"`
}

// verifyID will verify an account ID
func verifyID(id string) error {
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

	// Allow overriding the data path
	var datapath = *data
	var override = os.Getenv("USER_DATA")
	if len(override) > 0 {
		datapath = override
	}
	datapath = path.Join(datapath, "rove.json")

	// Create the path if needed
	path := filepath.Dir(datapath)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	} else {
		// Read the file
		_, err = os.Stat(datapath)
		if !os.IsNotExist(err) {
			if b, err := ioutil.ReadFile(datapath); err != nil {
				return fmt.Errorf("failed to read file %s error: %s", datapath, err)

			} else if len(b) == 0 {
				return fmt.Errorf("file %s was empty, assumin fresh data", datapath)

			} else if err := json.Unmarshal(b, &config); err != nil {
				return fmt.Errorf("failed to unmarshal file %s error: %s", datapath, err)

			}
		}
	}

	// Early bails
	switch command {
	case "version":
		fmt.Println(version.Version)
		return nil
	case "config":
		fmt.Printf("host: %s\taccount: %s\n", config.Host, config.Accounts[config.Host])
		return nil
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
		if len(*name) == 0 {
			return fmt.Errorf("must set name with -name")
		}
		d := rove.RegisterRequest{
			Name: *name,
		}
		_, err := client.Register(ctx, &d)
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Registered account with id: %s\n", *name)
			config.Accounts[config.Host] = *name
		}

	case "move":
		d := rove.CommandsRequest{
			Account: config.Accounts[config.Host],
			Commands: []*rove.Command{
				{
					Command: game.CommandMove,
					Bearing: *bearing,
				},
			},
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
			fmt.Printf("attributes: %+v\n", response)
		}

	default:
		// Print the usage
		fmt.Fprintf(os.Stderr, "Error: unknown command %s\n", command)
		printUsage()
		os.Exit(1)
	}

	// Save out the persistent file
	if b, err := json.MarshalIndent(config, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal data error: %s", err)
	} else if err := ioutil.WriteFile(datapath, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file %s error: %s", datapath, err)
	}

	return nil
}

// Simple main
func main() {
	flag.Usage = printUsage

	// Bail without any args
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	flag.CommandLine.Parse(os.Args[2:])

	// Run the inner main
	if err := InnerMain(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
