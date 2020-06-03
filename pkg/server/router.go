package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/version"
)

// Route defines the information for a single path->function route
type Route struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

// NewRouter sets up the server mux
func (s *Server) SetUpRouter() {

	// Array of all our routes
	var routes = []Route{
		{
			path:    "/status",
			handler: s.HandleStatus,
		},
		{
			path:    "/register",
			handler: s.HandleRegister,
		},
		{
			path:    "/spawn",
			handler: s.HandleSpawn,
		},
	}

	// Set up the handlers
	for _, route := range routes {
		s.router.HandleFunc(route.path, route.handler)
	}
}

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready   bool   `json:"ready"`
	Version string `json:"version"`
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
		Ready:   true,
		Version: version.Version,
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}

// BasicResponse describes the minimum dataset for a response
type BasicResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// BasicAccountData describes the data to be sent for an account specific post
type BasicAccountData struct {
	Id string `json:"id"`
}

// RegisterData describes the data to send when registering
type RegisterData struct {
	Name string `json:"name"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	BasicResponse

	Id string `json:"id"`
}

// HandleRegister handles HTTP requests to the /register endpoint
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

	// Set up the response
	var response = RegisterResponse{
		BasicResponse: BasicResponse{
			Success: false,
		},
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
		fmt.Printf("Failed to decode json: %s\n", err)

		response.Error = err.Error()
	} else if len(data.Name) == 0 {
		response.Error = "Cannot register empty name"
	} else {
		// log the data sent
		fmt.Printf("\tdata: %+v\n", data)

		// Register the account with the server
		acc, err := s.accountant.RegisterAccount(data.Name)

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

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}

// SpawnData is the data to be sent for the spawn command
type SpawnData struct {
	BasicAccountData
}

// SpawnResponse is the data to respond with on a spawn command
type SpawnResponse struct {
	BasicResponse

	Position game.Vector `json:"position"`
}

// HandleSpawn will spawn the player entity for the associated account
func (s *Server) HandleSpawn(w http.ResponseWriter, r *http.Request) {
	// Verify we're hit with a get request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

	// Set up the response
	var response = SpawnResponse{
		BasicResponse: BasicResponse{
			Success: false,
		},
	}

	// Pull out the incoming info
	var data SpawnData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(data.Id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(data.Id); err != nil {
		response.Error = "Provided account ID was invalid"

	} else {
		// log the data sent
		fmt.Printf("\tspawn data: %v\n", data)

		// Create a new instance
		inst := uuid.New()
		s.world.Spawn(inst)
		if pos, err := s.world.GetPosition(inst); err != nil {
			response.Error = fmt.Sprint("No position found for created instance")

		} else {
			if err := s.accountant.AssignPrimary(id, inst); err != nil {
				response.Error = err.Error()

				// Try and clear up the instance
				if err := s.world.DestroyInstance(inst); err != nil {
					fmt.Printf("Failed to destroy instance after failed primary assign: %s", err)
				}

			} else {
				// Reply with valid data
				response.Success = true
				response.Position = pos
			}
		}
	}

	// Be a good citizen and set the header for the return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)
}
