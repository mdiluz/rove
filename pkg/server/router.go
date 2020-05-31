package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter sets up the server mux
func (s *Server) SetUpRouter() {
	s.router = mux.NewRouter().StrictSlash(true)

	// Set up the handlers
	s.router.HandleFunc("/status", s.HandleStatus)
	s.router.HandleFunc("/register", s.HandleRegister)
}

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready bool `json:"ready"`
}

// HandleStatus handles HTTP requests to the /status endpoint
func (s *Server) HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\t%s", r.Method, r.RequestURI)

	// Verify we're hit with a get request
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var response = StatusResponse{
		Ready: true,
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}

// RegisterData describes the data to send when registering
type RegisterData struct {
	Name string `json:"id"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	Id string `json:"id"`

	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// HandleRegister handles HTTP requests to the /register endpoint
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\t%s", r.Method, r.RequestURI)

	// Verify we're hit with a get request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Pull out the registration info
	var data RegisterData
	json.NewDecoder(r.Body).Decode(&data)

	// Register the account with the server
	acc := Account{Name: data.Name}
	acc, err := s.accountant.RegisterAccount(acc)

	// Set up the response
	var response = RegisterResponse{
		Success: false,
	}

	// If we didn't fail, respond with the account ID string
	if err == nil {
		response.Success = true
		response.Id = acc.id.String()
	} else {
		response.Error = err.Error()
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}
