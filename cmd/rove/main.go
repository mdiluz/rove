package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mdiluz/rove/pkg/atlas"
	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/pkg/objects"
	"github.com/mdiluz/rove/pkg/roveapi"
	"github.com/mdiluz/rove/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var home = os.Getenv("HOME")
var defaultDataPath = path.Join(home, ".local/share/")

// Command usage
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: rove COMMAND [ARGS...]\n")
	fmt.Fprintln(os.Stderr, "\nCommands")
	fmt.Fprintln(os.Stderr, "\tserver-status              prints the server status")
	fmt.Fprintln(os.Stderr, "\tregister NAME              registers an account and stores it (use with -name)")
	fmt.Fprintln(os.Stderr, "\tcommand COMMAND [VAL...]   issue commands to rover, accepts multiple, see below")
	fmt.Fprintln(os.Stderr, "\tradar                      gathers radar data for the current rover")
	fmt.Fprintln(os.Stderr, "\tstatus                     gets status info for current rover")
	fmt.Fprintln(os.Stderr, "\tconfig [HOST]              outputs the local config info, optionally sets host")
	fmt.Fprintln(os.Stderr, "\thelp                       outputs this usage information")
	fmt.Fprintln(os.Stderr, "\tversion                    outputs version info")
	fmt.Fprintln(os.Stderr, "\nRover commands:")
	fmt.Fprintln(os.Stderr, "\tmove BEARING               moves the rover in the chosen direction")
	fmt.Fprintln(os.Stderr, "\tstash                      stores the object at the rover location in the inventory")
	fmt.Fprintln(os.Stderr, "\trepair                     uses an inventory object to repair the rover")
	fmt.Fprintln(os.Stderr, "\trecharge                   wait a tick to recharge the rover")
	fmt.Fprintln(os.Stderr, "\tbroadcast MSG              broadcast a simple ASCII triplet to nearby rovers")
	fmt.Fprintln(os.Stderr, "\nEnvironment")
	fmt.Fprintln(os.Stderr, "\tROVE_USER_DATA             path to user data, defaults to "+defaultDataPath)
}

const gRPCport = 9090

// Account stores data for an account
type Account struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

// Config is used to store internal data
type Config struct {
	Host    string  `json:"host,omitempty"`
	Account Account `json:"account,omitempty"`
}

// ConfigPath returns the configuration path
func ConfigPath() string {
	// Allow overriding the data path
	var datapath = defaultDataPath
	var override = os.Getenv("ROVE_USER_DATA")
	if len(override) > 0 {
		datapath = override
	}
	datapath = path.Join(datapath, "roveapi.json")

	return datapath
}

// LoadConfig loads the config from a chosen path
func LoadConfig() (config Config, err error) {

	datapath := ConfigPath()

	// Create the path if needed
	path := filepath.Dir(datapath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return Config{}, fmt.Errorf("Failed to create data path %s: %s", path, err)
		}
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

// checkAccount will verify an account ID
func checkAccount(a Account) error {
	if len(a.Name) == 0 {
		return fmt.Errorf("no account ID set, must register first")
	} else if len(a.Secret) == 0 {
		return fmt.Errorf("empty account secret, must register first")
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
		fmt.Printf("host: %s\taccount: %s\n", config.Host, config.Account)
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
	var client = roveapi.NewRoveClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Handle all the commands
	switch command {
	case "server-status":
		response, err := client.ServerStatus(ctx, &roveapi.ServerStatusRequest{})
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Ready: %t\n", response.Ready)
			fmt.Printf("Version: %s\n", response.Version)
			fmt.Printf("Tick Rate: %d\n", response.TickRate)
			fmt.Printf("Current Tick: %d\n", response.CurrentTick)
			fmt.Printf("Next Tick: %s\n", response.NextTick)
		}

	case "register":
		if len(args) == 0 || len(args[0]) == 0 {
			return fmt.Errorf("must pass name to 'register'")
		}

		resp, err := client.Register(ctx, &roveapi.RegisterRequest{
			Name: args[0],
		})
		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Registered account with id: %s\n", resp.Account.Name)
			config.Account.Name = resp.Account.Name
			config.Account.Secret = resp.Account.Secret
		}

	case "command":
		if err := checkAccount(config.Account); err != nil {
			return err
		} else if len(args) == 0 {
			return fmt.Errorf("must pass commands to 'commands'")
		}

		// Iterate through each command
		var commands []*roveapi.Command
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "move":
				i++
				if len(args) == i {
					return fmt.Errorf("move command must be passed bearing")
				} else if _, err := maths.FromString(args[i]); err != nil {
					return err
				}
				commands = append(commands,
					&roveapi.Command{
						Command: roveapi.CommandType_move,
						Data:    &roveapi.Command_Bearing{Bearing: args[i]},
					},
				)
			case "broadcast":
				i++
				if len(args) == i {
					return fmt.Errorf("broadcast command must be passed an ASCII triplet")
				} else if len(args[i]) > 3 {
					return fmt.Errorf("broadcast command must be given ASCII triplet of 3 or less: %s", args[i])
				}
				commands = append(commands,
					&roveapi.Command{
						Command: roveapi.CommandType_broadcast,
						Data:    &roveapi.Command_Message{Message: []byte(args[i])},
					},
				)
			default:
				// By default just use the command literally
				commands = append(commands,
					&roveapi.Command{
						Command: roveapi.CommandType(roveapi.CommandType_value[args[i]]),
					},
				)
			}
		}

		_, err := client.Command(ctx, &roveapi.CommandRequest{
			Account: &roveapi.Account{
				Name:   config.Account.Name,
				Secret: config.Account.Secret,
			},
			Commands: commands,
		})

		switch {
		case err != nil:
			return err

		default:
			fmt.Printf("Request succeeded\n")
		}

	case "radar":
		if err := checkAccount(config.Account); err != nil {
			return err
		}

		response, err := client.Radar(ctx, &roveapi.RadarRequest{
			Account: &roveapi.Account{
				Name:   config.Account.Name,
				Secret: config.Account.Secret,
			},
		})

		switch {
		case err != nil:
			return err

		default:

			// Print out the radar
			num := int(math.Sqrt(float64(len(response.Tiles))))
			for j := num - 1; j >= 0; j-- {
				for i := 0; i < num; i++ {
					t := response.Tiles[i+num*j]
					o := response.Objects[i+num*j]
					if o != byte(objects.None) {
						fmt.Printf("%c", o)
					} else if t != byte(atlas.TileNone) {
						fmt.Printf("%c", t)
					} else {
						fmt.Printf(" ")
					}

				}
				fmt.Print("\n")
			}
		}

	case "status":
		if err := checkAccount(config.Account); err != nil {
			return err
		}

		response, err := client.Status(ctx, &roveapi.StatusRequest{
			Account: &roveapi.Account{
				Name:   config.Account.Name,
				Secret: config.Account.Secret,
			},
		})

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
