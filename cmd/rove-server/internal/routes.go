package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
	"google.golang.org/grpc"
)

// Handler describes a function that handles any incoming request and can respond
type Handler func(*Server, map[string]string, io.ReadCloser) (interface{}, error)

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
func HandleStatus(s *Server, vars map[string]string, b io.ReadCloser) (interface{}, error) {

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
func HandleRegister(s *Server, vars map[string]string, b io.ReadCloser) (interface{}, error) {
	var response = rove.RegisterResponse{}

	// Decode the registration info, verify it and register the account
	var data rove.RegisterData
	err := json.NewDecoder(b).Decode(&data)
	if err != nil {
		log.Printf("Failed to decode json: %s\n", err)
		return BadRequestError{Error: err.Error()}, nil

	} else if len(data.Name) == 0 {
		return BadRequestError{Error: "cannot register empty name"}, nil

	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	reg := accounts.RegisterInfo{Name: data.Name}
	if acc, err := s.accountant.Register(ctx, &reg, grpc.WaitForReady(true)); err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)

	} else if !acc.Success {
		return BadRequestError{Error: acc.Error}, nil

	} else if _, _, err := s.SpawnRoverForAccount(data.Name); err != nil {
		return nil, fmt.Errorf("failed to spawn rover for account: %s", err)

	} else if err := s.SaveWorld(); err != nil {
		return nil, fmt.Errorf("internal server error when saving world: %s", err)

	}

	log.Printf("register response:%+v\n", response)
	return response, nil
}

// HandleSpawn will spawn the player entity for the associated account
func HandleCommand(s *Server, vars map[string]string, b io.ReadCloser) (interface{}, error) {
	var response = rove.CommandResponse{}

	id := vars["account"]

	// Decode the commands, verify them and the account, and execute the commands
	var data rove.CommandData
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		log.Printf("Failed to decode json: %s\n", err)
		return BadRequestError{Error: err.Error()}, nil

	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		return BadRequestError{Error: "no account ID provided"}, nil

	} else if resp, err := s.accountant.GetValue(ctx, &key); err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)

	} else if !resp.Success {
		return BadRequestError{Error: resp.Error}, nil

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)

	} else if err := s.world.Enqueue(id, data.Commands...); err != nil {
		return BadRequestError{Error: err.Error()}, nil

	}

	log.Printf("command response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}

// HandleRadar handles the radar request
func HandleRadar(s *Server, vars map[string]string, b io.ReadCloser) (interface{}, error) {
	var response = rove.RadarResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	id := vars["account"]
	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		return BadRequestError{Error: "no account ID provided"}, nil

	} else if resp, err := s.accountant.GetValue(ctx, &key); err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)

	} else if !resp.Success {
		return BadRequestError{Error: resp.Error}, nil

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)

	} else if attrib, err := s.world.RoverAttributes(id); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else if radar, err := s.world.RadarFromRover(id); err != nil {
		return nil, fmt.Errorf("error getting radar from rover: %s", err)

	} else {
		response.Tiles = radar
		response.Range = attrib.Range
	}

	log.Printf("radar response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}

// HandleRover handles the rover request
func HandleRover(s *Server, vars map[string]string, b io.ReadCloser) (interface{}, error) {
	var response = rove.RoverResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	id := vars["account"]
	key := accounts.DataKey{Account: id, Key: "rover"}
	if len(id) == 0 {
		return BadRequestError{Error: "no account ID provided"}, nil

	} else if resp, err := s.accountant.GetValue(ctx, &key); err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)

	} else if !resp.Success {
		return BadRequestError{Error: resp.Error}, nil

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)

	} else if attrib, err := s.world.RoverAttributes(id); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else {
		response.Attributes = attrib
	}

	log.Printf("rover response \taccount:%s\tresponse:%+v\n", id, response)
	return response, nil
}
