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
	server := server.NewServer(*port)

	fmt.Println("Initialising...")

	server.Initialise()

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("SIGTERM recieved, exiting...")
		os.Exit(0)
	}()

	fmt.Println("Initialised")

	server.Run()
}
