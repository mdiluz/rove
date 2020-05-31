package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mdiluz/rove/pkg/accounts"
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
	fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

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
	fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

	// Set up the response
	var response = RegisterResponse{
		Success: false,
	}

	// Verify we're hit with a get request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Pull out the registration info
	var data RegisterData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to decode json: %s", err)

		response.Error = err.Error()
	} else {
		// log the data sent
		fmt.Printf("\t%v\n", data)

		// Register the account with the server
		acc := accounts.Account{Name: data.Name}
		acc, err := s.accountant.RegisterAccount(acc)

		// If we didn't fail, respond with the account ID string
		if err == nil {
			response.Success = true
			response.Id = acc.Id.String()
		} else {
			response.Error = err.Error()
		}
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}
