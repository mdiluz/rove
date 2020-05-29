package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	fmt.Println("Initialising...")

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("SIGTERM recieved, exiting...")
		os.Exit(0)
	}()

	// Create a new router
	router := NewRouter()

	fmt.Println("Initialised")

	// Listen and serve the http requests
	fmt.Println("Serving HTTP")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
