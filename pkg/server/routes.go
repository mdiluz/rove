package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
)

// Handler describes a function that handles any incoming request and can respond
type Handler func(*Server, map[string]string, io.ReadCloser, io.Writer) (interface{}, error)

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
		path:    "/{account}/spawn",
		method:  http.MethodPost,
		handler: HandleSpawn,
	},
	{
		path:    "/{account}/command",
		method:  http.MethodPost,
		handler: HandleCommand,
	},
	{
		path:    "/{account}/radar",
		method:  http.MethodGet,
		handler: HandleRadar,
	},
	{
		path:    "/{account}/rover",
		method:  http.MethodGet,
		handler: HandleRover,
	},
}

// HandleStatus handles the /status request
func HandleStatus(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {

	// Simply return the current server status
	return rove.StatusResponse{
		Ready:   true,
		Version: version.Version,
	}, nil
}

// HandleRegister handles /register endpoint
func HandleRegister(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.RegisterResponse{
		Success: false,
	}

	// Decode the registration info, verify it and register the account
	var data rove.RegisterData
	err := json.NewDecoder(b).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(data.Name) == 0 {
		response.Error = "Cannot register empty name"

	} else if acc, err := s.accountant.RegisterAccount(data.Name); err != nil {
		response.Error = err.Error()

	} else if err := s.SaveAll(); err != nil {
		response.Error = fmt.Sprintf("Internal server error when saving accounts: %s", err)

	} else {
		// Save out the new accounts
		fmt.Printf("New account registered\tname:%s\tid:%s", acc.Name, acc.Id)

		response.Id = acc.Id.String()
		response.Success = true
	}

	return response, nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleSpawn(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.SpawnResponse{
		Success: false,
	}

	id := vars["account"]

	// Decode the spawn info, verify it and spawn the rover for this account
	var data rove.SpawnData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(id); err != nil {
		response.Error = "Provided account ID was invalid"

	} else if pos, rover, err := s.SpawnRoverForAccount(id); err != nil {
		response.Error = err.Error()

	} else {
		fmt.Printf("New rover spawned\taccount:%s\trover:%s\tpos:%+v", id, rover, pos)

		response.Success = true
		response.Position = pos
	}

	return response, nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleCommand(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.CommandResponse{
		Success: false,
	}

	id := vars["account"]

	// Decode the commands, verify them and the account, and execute the commands
	var data rove.CommandData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		fmt.Printf("Failed to decode json: %s\n", err)
		response.Error = err.Error()

	} else if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(id); err != nil {
		response.Error = fmt.Sprintf("Provided account ID was invalid: %s", err)

	} else if inst, err := s.accountant.GetRover(id); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if err := s.world.Enqueue(inst, data.Commands...); err != nil {
		response.Error = fmt.Sprintf("Failed to execute commands: %s", err)

	} else {
		fmt.Printf("Queued commands\taccount:%s\tcommands:%+v", id, data.Commands)
		response.Success = true
	}

	return response, nil
}

// HandleRadar handles the radar request
func HandleRadar(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.RadarResponse{
		Success: false,
	}

	id := vars["account"]
	if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(id); err != nil {
		response.Error = fmt.Sprintf("Provided account ID was invalid: %s", err)

	} else if inst, err := s.accountant.GetRover(id); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if radar, err := s.world.RadarFromRover(inst); err != nil {
		response.Error = fmt.Sprintf("Error getting radar from rover: %s", err)

	} else {
		fmt.Printf("Responded with radar\taccount:%s\tradar:%+v", id, radar)
		response.Rovers = radar.Rovers
		response.Success = true
	}

	return response, nil
}

// HandleRover handles the rover request
func HandleRover(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.RoverResponse{
		Success: false,
	}

	id := vars["account"]
	if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if id, err := uuid.Parse(id); err != nil {
		response.Error = fmt.Sprintf("Provided account ID was invalid: %s", err)

	} else if inst, err := s.accountant.GetRover(id); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if pos, err := s.world.RoverPosition(inst); err != nil {
		response.Error = fmt.Sprintf("Error getting radar from rover: %s", err)

	} else {
		fmt.Printf("Responded with rover\taccount:%s\trover:%+v", id, pos)
		response.Position = pos
		response.Success = true
	}

	return response, nil
}
