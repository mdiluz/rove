package main

import (
	"log"
	"net/http"
)

func main() {

	// Create a new router
	router := NewRouter()

	// Listen and serve the http requests
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
