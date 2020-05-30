package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/rove"
)

// NewRouter sets up the server mux
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)

	// Set up the handlers
	router.HandleFunc("/status", HandleStatus)

	return
}

// HandleStatus handles HTTP requests to the /status endpoint
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	var status = rove.ServerStatus{
		Ready: true,
	}

	json.NewEncoder(w).Encode(status)
}
