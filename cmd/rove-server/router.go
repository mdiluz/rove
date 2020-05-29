package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() (router *mux.Router) {
	// Create a new router
	router = mux.NewRouter().StrictSlash(true)

	// Set up the handlers
	router.HandleFunc("/status", HandleStatus)

	return
}

// RouterStatus is a struct that contains information on the status of the server
type RouterStatus struct {
	Ready bool `json:"ready"`
}

// HandleStatus handles HTTP requests to the /status endpoint
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	var status = RouterStatus{
		Ready: true,
	}

	json.NewEncoder(w).Encode(status)
}
