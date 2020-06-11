package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mdiluz/rove/cmd/rove-server/internal"
	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/version"
)

var ver = flag.Bool("version", false, "Display version number")
var quit = flag.Int("quit", 0, "Quit after n seconds, useful for testing")

// Address to host the server on, automatically selected if empty
var address = os.Getenv("HOST_ADDRESS")

// Path for persistent storage
var data = os.Getenv("DATA_PATH")

// The tick rate of the server in seconds
var tick = os.Getenv("TICK_RATE")

func InnerMain() {
	flag.Parse()

	// Print the version if requested
	if *ver {
		fmt.Println(version.Version)
		return
	}

	fmt.Printf("Initialising version %s...\n", version.Version)

	// Set the persistence path
	persistence.SetPath(data)

	// Convert the tick rate
	tickRate := 5
	if len(tick) > 0 {
		var err error
		tickRate, err = strconv.Atoi(tick)
		if err != nil {
			log.Fatalf("TICK_RATE not set to valid int: %s", err)
		}
	}

	// Create the server data
	s := internal.NewServer(
		internal.OptionAddress(address),
		internal.OptionPersistentData(),
		internal.OptionTick(tickRate))

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
	InnerMain()
}
