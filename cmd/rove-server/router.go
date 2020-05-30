package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/rove"
)

// NewRouter sets up the server mux
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)

	// Set up the handlers
	router.HandleFunc("/status", HandleStatus)
	router.HandleFunc("/register", HandleRegister)

	return
}

// HandleStatus handles HTTP requests to the /status endpoint
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\t%s", r.Method, r.RequestURI)

	var response = rove.StatusResponse{
		Ready: true,
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}

// HandleRegister handles HTTP requests to the /register endpoint
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\t%s", r.Method, r.RequestURI)

	// TODO: Add this user to the server
	player := NewPlayer()
	var response = rove.RegisterResponse{
		Success: true,
		Id:      player.id.String(),
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}
