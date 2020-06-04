package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/version"
)

// Handler describes a function that handles any incoming request and can respond
type Handler func(*Server, io.ReadCloser, io.Writer) error

// Route defines the information for a single path->function route
type Route struct {
	path    string
	method  string
	handler Handler
}

// Routes is an array of all the Routes
var Routes = []Route{
	{
		path:    "/status",
		method:  http.MethodGet,
		handler: HandleStatus,
	},
	{
		path:    "/register",
		method:  http.MethodPost,
		handler: HandleRegister,
	},
	{
		path:    "/spawn",
		method:  http.MethodPost,
		handler: HandleSpawn,
	},
	{
		path:    "/commands",
		method:  http.MethodPost,
		handler: HandleCommands,
	},
	{
		path:    "/view",
		method:  http.MethodPost,
		handler: HandleView,
	},
}

// HandleStatus handles the /status request
func HandleStatus(s *Server, b io.ReadCloser, w io.Writer) error {

	// Simply encode the current status
	var response = StatusResponse{
		Ready:   true,
		Version: version.Version,
	}

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// HandleRegister handles /register endpoint
func HandleRegister(s *Server, b io.ReadCloser, w io.Writer) error {

	// Set up the response
	var response = RegisterResponse{
		Success: false,
	}

	// Pull out the registration info
	var data RegisterData
	err := json.NewDecoder(b).Decode(&data)
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

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleSpawn(s *Server, b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = SpawnResponse{
		Success: false,
	}

	// Pull out the incoming info
	var data SpawnData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(data.Id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(data.Id); err != nil {
		response.Error = "Provided account ID was invalid"

	} else {
		// log the data sent
		fmt.Printf("\tspawn data: %v\n", data)

		// Create a new rover
		if pos, _, err := s.SpawnRoverForAccount(id); err != nil {
			response.Error = err.Error()
		} else {
			response.Success = true
			response.X = pos.X
			response.Y = pos.Y
		}
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleCommands(s *Server, b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = CommandsResponse{
		Success: false,
	}

	// Pull out the incoming info
	var data CommandsData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(data.Id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(data.Id); err != nil {
		response.Error = fmt.Sprintf("Provided account ID was invalid: %s", err)

	} else if inst, err := s.accountant.GetRover(id); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else {
		// log the data sent
		fmt.Printf("\tcommands data: %v\n", data)

		// Iterate through the commands to generate all game commands
		var cmds []game.Command
		for _, c := range data.Commands {
			switch c.Command {
			case CommandMove:
				cmds = append(cmds, s.world.CommandMove(inst, c.Bearing, c.Duration))
			}
		}

		// Execute the commands
		if err := s.world.Execute(cmds...); err != nil {
			response.Error = fmt.Sprintf("Failed to execute commands: %s", err)
		} else {
			response.Success = true
		}
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// HandleView handles the view request
func HandleView(s *Server, b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = ViewResponse{
		Success: false,
	}

	// Pull out the incoming info
	var data CommandsData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(data.Id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(data.Id); err != nil {
		response.Error = fmt.Sprintf("Provided account ID was invalid: %s", err)

	} else {
		// log the data sent
		fmt.Printf("\tcommands data: %v\n", data)

		// TODO: Query the view for this account
		fmt.Println(id)
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}
