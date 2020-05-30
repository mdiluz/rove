package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdiluz/rove/pkg/rovegame"
)

var port = flag.Int("port", 8080, "The port to host on")

func main() {

	fmt.Println("Initialising...")

	// Set up the world
	world := rovegame.NewWorld()
	fmt.Printf("World created\n\t%+v\n", world)

	// Create a new router
	router := NewRouter()
	fmt.Printf("Router Created\n")

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("SIGTERM recieved, exiting...")
		os.Exit(0)
	}()

	fmt.Println("Initialised")

	// Listen and serve the http requests
	fmt.Println("Serving HTTP")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Fatal(err)
	}
}
