package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/server"
	"github.com/mdiluz/rove/pkg/version"
)

var ver = flag.Bool("version", false, "Display version number")
var port = flag.String("address", ":8080", "The address to host on")
var data = flag.String("data", os.TempDir(), "Directory to store persistant data")
var quit = flag.Int("quit", 0, "Quit after n seconds, useful for testing")

func main() {
	flag.Parse()

	// Print the version if requested
	if *ver {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	fmt.Printf("Initialising version %s...\n", version.Version)

	// Set the persistence path
	persistence.SetPath(*data)

	// Create the server data
	s := server.NewServer(
		server.OptionAddress(*port),
		server.OptionPersistentData())

	// Initialise the server
	if err := s.Initialise(); err != nil {
		panic(err)
	}

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("SIGTERM recieved, exiting...")
		if err := s.Close(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	// Quit after a time if requested
	if *quit != 0 {
		go func() {
			time.Sleep(time.Duration(*quit) * time.Second)
			if err := s.Close(); err != nil {
				panic(err)
			}
			os.Exit(0)
		}()
	}

	// Run the server
	fmt.Printf("Serving HTTP on %s\n", s.Addr())
	s.Run()

	// Close the server
	if err := s.Close(); err != nil {
		panic(err)
	}
}
