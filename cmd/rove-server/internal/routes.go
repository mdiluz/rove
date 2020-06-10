package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/accounts"
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
	response := rove.StatusResponse{
		Ready:   true,
		Version: version.Version,
		Tick:    s.tick,
	}

	// If there's a schedule, respond with it
	if len(s.schedule.Entries()) > 0 {
		response.NextTick = s.schedule.Entries()[0].Next.Format("15:04:05")
	}

	return response, nil
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

	}
	reg := accounts.RegisterInfo{Name: data.Name}
	if acc, err := s.accountant.Register(context.Background(), &reg); err != nil {
		response.Error = err.Error()

	} else if !acc.Success {
		response.Error = acc.Error

	} else if _, _, err := s.SpawnRoverForAccount(data.Name); err != nil {
		response.Error = err.Error()

	} else if err := s.SaveWorld(); err != nil {
		response.Error = fmt.Sprintf("Internal server error when saving world: %s", err)

	} else {
		// Save out the new accounts
		response.Success = true
	}

	fmt.Printf("register response:%+v\n", response)
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

	}

	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if resp, err := s.accountant.GetValue(context.Background(), &key); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if !resp.Success {
		response.Error = resp.Error

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		response.Error = fmt.Sprintf("Account had invalid rover id: %s", err)

	} else if err := s.world.Enqueue(id, data.Commands...); err != nil {
		response.Error = fmt.Sprintf("Failed to execute commands: %s", err)

	} else {
		response.Success = true
	}

	fmt.Printf("command response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}

// HandleRadar handles the radar request
func HandleRadar(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.RadarResponse{
		Success: false,
	}

	id := vars["account"]
	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if resp, err := s.accountant.GetValue(context.Background(), &key); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if !resp.Success {
		response.Error = resp.Error

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		response.Error = fmt.Sprintf("Account had invalid rover id: %s", err)

	} else if attrib, err := s.world.RoverAttributes(id); err != nil {
		response.Error = fmt.Sprintf("Error getting rover attributes: %s", err)

	} else if radar, err := s.world.RadarFromRover(id); err != nil {
		response.Error = fmt.Sprintf("Error getting radar from rover: %s", err)

	} else {
		response.Tiles = radar
		response.Range = attrib.Range
		response.Success = true
	}

	fmt.Printf("radar response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}

// HandleRover handles the rover request
func HandleRover(s *Server, vars map[string]string, b io.ReadCloser, w io.Writer) (interface{}, error) {
	var response = rove.RoverResponse{
		Success: false,
	}

	id := vars["account"]
	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		response.Error = "No account ID provided"

	} else if resp, err := s.accountant.GetValue(context.Background(), &key); err != nil {
		response.Error = fmt.Sprintf("Provided account has no rover: %s", err)

	} else if !resp.Success {
		response.Error = resp.Error

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		response.Error = fmt.Sprintf("Account had invalid rover id: %s", err)

	} else if attribs, err := s.world.RoverAttributes(id); err != nil {
		response.Error = fmt.Sprintf("Error getting radar from rover: %s", err)

	} else {
		response.Attributes = attribs
		response.Success = true
	}

	fmt.Printf("rover response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}
