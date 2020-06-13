package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Command usage
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [OPTIONS]...\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "\tstatus  \tprints the server status")
	fmt.Fprintln(os.Stderr, "\tregister\tregisters an account and stores it (use with -name)")
	fmt.Fprintln(os.Stderr, "\tmove    \tissues move command to rover")
	fmt.Fprintln(os.Stderr, "\tradar   \tgathers radar data for the current rover")
	fmt.Fprintln(os.Stderr, "\trover   \tgets data for current rover")
	fmt.Fprintln(os.Stderr, "\tconfig  \toutputs the local config info")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}

var home = os.Getenv("HOME")
var filepath = path.Join(home, ".local/share/rove.json")

const gRPCport = 9090

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
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Host, gRPCport), grpc.WithInsecure())
	if err != nil {
		return err
	}
	var client = rove.NewRoveClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Grab the account
	var account = config.Accounts[config.Host]

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
			Account: account,
			Commands: []*rove.Command{
				{
					Command:  game.CommandMove,
					Duration: int32(*duration),
					Bearing:  *bearing,
				},
			},
		}

		if err := verifyId(account); err != nil {
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
		dat := rove.RadarRequest{Account: account}
		if err := verifyId(account); err != nil {
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
		req := rove.RoverRequest{Account: account}
		if err := verifyId(account); err != nil {
			return err
		}
		response, err := client.Rover(ctx, &req)

		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("attributes: %+v\n", response)
		}
	case "config":
		fmt.Printf("host: %s\taccount: %s\n", config.Host, account)

	default:
		// Print the usage
		fmt.Fprintf(os.Stderr, "Error: unknown command %s\n", command)
		Usage()
		os.Exit(1)
	}

	// Save out the persistent file
	if b, err := json.MarshalIndent(config, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal data error: %s", err)
	} else if err := ioutil.WriteFile(*data, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file %s error: %s", *data, err)
	}

	return nil
}

// Simple main
func main() {
	flag.Usage = Usage

	// Bail without any args
	if len(os.Args) == 1 {
		Usage()
		os.Exit(1)
	}

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
