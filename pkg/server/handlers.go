package server

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/version"
)

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready   bool   `json:"ready"`
	Version string `json:"version"`
}

// HandleStatus handles the /status request
func (s *Server) HandleStatus(b io.ReadCloser, w io.Writer) error {

	// Simply encode the current status
	var response = StatusResponse{
		Ready:   true,
		Version: version.Version,
	}

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
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

// HandleRegister handles /register endpoint
func (s *Server) HandleRegister(b io.ReadCloser, w io.Writer) error {

	// Set up the response
	var response = RegisterResponse{
		BasicResponse: BasicResponse{
			Success: false,
		},
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
func (s *Server) HandleSpawn(b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = SpawnResponse{
		BasicResponse: BasicResponse{
			Success: false,
		},
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

		// Create a new instance
		if pos, _, err := s.SpawnPrimary(id); err != nil {
			response.Error = err.Error()
		} else {
			response.Success = true
			response.Position = pos
		}
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

const (
	// CommandMove describes a single move command
	CommandMove = "move"
)

// Command describes a single command to execute
// it contains the type, and then any members used for each command type
type Command struct {
	// Command is the main command string
	Command string `json:"command"`

	// Used for CommandMove
	Vector game.Vector `json:"vector"`
}

// CommandsData is a set of commands to execute in order
type CommandsData struct {
	BasicAccountData
	Commands []Command `json:"commands"`
}

// HandleSpawn will spawn the player entity for the associated account
func (s *Server) HandleCommands(b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = BasicResponse{
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

	} else if inst, err := s.accountant.GetPrimary(id); err != nil {
		response.Error = fmt.Sprintf("Provided account has no primary: %s", err)

	} else {
		// log the data sent
		fmt.Printf("\tcommands data: %v\n", data)

		// Iterate through the commands to generate all game commands
		var cmds []game.Command
		for _, c := range data.Commands {
			switch c.Command {
			case CommandMove:
				cmds = append(cmds, s.world.CommandMove(inst, c.Vector))
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

// ViewData describes the input data to request an accounts current view
type ViewData struct {
	BasicAccountData
}

// ViewResponse describes the response to a /view call
type ViewResponse struct {
	BasicResponse
}

// HandleView handles the view request
func (s *Server) HandleView(b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = ViewResponse{
		BasicResponse: BasicResponse{
			Success: false,
		},
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
