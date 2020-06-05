package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/server"
	"github.com/mdiluz/rove/pkg/version"
)

var ver = flag.Bool("version", false, "Display version number")
var port = flag.Int("port", 8080, "The port to host on")
var data = flag.String("data", os.TempDir(), "Directory to store persistant data")

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
		server.OptionPort(*port),
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

	fmt.Println("Running...")

	// Run the server
	s.Run()

	// Close the server
	if err := s.Close(); err != nil {
		panic(err)
	}
}
