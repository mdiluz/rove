package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mdiluz/rove/cmd/rove-server/internal"
	"github.com/mdiluz/rove/pkg/persistence"
	"github.com/mdiluz/rove/pkg/version"
)

var ver = flag.Bool("version", false, "Display version number")

// Path for persistent storage
var data = os.Getenv("DATA_PATH")

// The tick rate of the server in seconds
var tick = os.Getenv("TICK_RATE")

func InnerMain() {
	flag.Parse()

	// Print the version if requested
	if *ver {
		log.Println(version.Version)
		return
	}

	// Address to host the server on
	var iport int
	var port = os.Getenv("PORT")
	if len(port) == 0 {
		iport = 9090
	} else {
		var err error
		iport, err = strconv.Atoi(port)
		if err != nil {
			log.Fatal("$PORT not valid int")
		}
	}
	log.Printf("Initialising version %s...\n", version.Version)

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
		internal.OptionAddress(fmt.Sprintf(":%d", iport)),
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
		log.Println("Quit requested, exiting...")
		if err := s.Stop(); err != nil {
			panic(err)
		}
	}()

	// Run the server
	s.Run()

	// Close the server
	if err := s.Close(); err != nil {
		panic(err)
	}
}

func main() {
	InnerMain()
}
