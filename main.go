package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdiluz/rove/pkg/server"
)

var port = flag.Int("port", 8080, "The port to host on")

func main() {
	s := server.NewServer(
		server.OptionPort(*port),
		server.OptionPersistentData())

	fmt.Println("Initialising...")
	if err := s.Initialise(); err != nil {
		panic(err)
	}

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("SIGTERM recieved, exiting...")
		s.Close()
		os.Exit(0)
	}()

	fmt.Println("Initialised")

	s.Run()

	if err := s.Close(); err != nil {
		panic(err)
	}
}
