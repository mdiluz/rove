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
		path:    "/radar",
		method:  http.MethodPost,
		handler: HandleRadar,
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

	} else if acc, err := s.accountant.RegisterAccount(data.Name); err != nil {
		response.Error = err.Error()

	} else {
		// log the data sent
		fmt.Printf("\tdata: %+v\n", data)

		// If we didn't fail, respond with the account ID string
		response.Id = acc.Id.String()
		response.Success = true
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

	} else if pos, _, err := s.SpawnRoverForAccount(id); err != nil {
		response.Error = err.Error()

	} else {

		// Create a new rover
		response.Success = true
		response.Position = pos

	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// ConvertCommands converts server commands to game commands
func (s *Server) ConvertCommands(commands []Command, inst uuid.UUID) ([]game.Command, error) {
	var cmds []game.Command
	for _, c := range commands {
		switch c.Command {
		case CommandMove:
			if bearing, err := game.DirectionFromString(c.Bearing); err != nil {
				return nil, err
			} else {
				cmds = append(cmds, s.world.CommandMove(inst, bearing, c.Duration))
			}
		}
	}
	return cmds, nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleCommands(s *Server, b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = CommandsResponse{
		Success: false,
	}

	// Go through each possible validation step for the data
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

	} else if cmds, err := s.ConvertCommands(data.Commands, inst); err != nil {
		response.Error = fmt.Sprintf("Couldn't convert commands: %s", err)

	} else if err := s.world.Execute(cmds...); err != nil {
		response.Error = fmt.Sprintf("Failed to execute commands: %s", err)

	} else {
		response.Success = true
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}

// HandleRadar handles the radar request
func HandleRadar(s *Server, b io.ReadCloser, w io.Writer) error {
	// Set up the response
	var response = RadarResponse{
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

	} else if radar, err := s.world.RadarFromRover(inst); err != nil {
		response.Error = fmt.Sprintf("Error getting radar from rover: %s", err)

	} else {
		// Fill in the response
		response.Rovers = radar.Rovers
		response.Success = true
	}

	// Log the response
	fmt.Printf("\tresponse: %+v\n", response)

	// Reply with the current status
	json.NewEncoder(w).Encode(response)

	return nil
}
