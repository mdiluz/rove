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
var quit = flag.Int("quit", 0, "Quit after n seconds, useful for testing")
var address = flag.String("address", "", "The address to host on, automatically selected if empty")
var data = flag.String("data", "", "Directory to store persistant data, no storage if empty")
var tick = flag.Int("tick", 5, "Number of minutes per server tick (0 for no tick)")

func InnerMain() {
	flag.Parse()

	// Print the version if requested
	if *ver {
		fmt.Println(version.Version)
		return
	}

	fmt.Printf("Initialising version %s...\n", version.Version)

	// Set the persistence path
	persistence.SetPath(*data)

	// Create the server data
	s := server.NewServer(
		server.OptionAddress(*address),
		server.OptionPersistentData(),
		server.OptionTick(*tick))

	// Initialise the server
	if err := s.Initialise(true); err != nil {
		panic(err)
	}

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Quit requested, exiting...")
		if err := s.Stop(); err != nil {
			panic(err)
		}
	}()

	// Quit after a time if requested
	if *quit != 0 {
		go func() {
			time.Sleep(time.Duration(*quit) * time.Second)
			syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
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

func main() {
	flag.Parse()
	InnerMain()
}
