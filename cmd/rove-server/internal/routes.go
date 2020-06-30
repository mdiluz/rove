package internal

import (
	"context"
	"fmt"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
)

// Status returns the status of the current server to a gRPC request
func (s *Server) Status(context.Context, *rove.StatusRequest) (*rove.StatusResponse, error) {
	response := &rove.StatusResponse{
		Ready:   true,
		Version: version.Version,
		Tick:    int32(s.tick),
	}

	// TODO: Verify the accountant is up and ready too

	// If there's a schedule, respond with it
	if len(s.schedule.Entries()) > 0 {
		response.NextTick = s.schedule.Entries()[0].Next.Format("15:04:05")
	}

	return response, nil
}

// Register registers a new account for a gRPC request
func (s *Server) Register(ctx context.Context, req *rove.RegisterRequest) (*rove.RegisterResponse, error) {
	if len(req.Name) == 0 {
		return nil, fmt.Errorf("empty account name")
	}

	if _, err := s.accountant.RegisterAccount(req.Name); err != nil {
		return nil, err

	} else if _, err := s.SpawnRoverForAccount(req.Name); err != nil {
		return nil, fmt.Errorf("failed to spawn rover for account: %s", err)

	} else if err := s.SaveWorld(); err != nil {
		return nil, fmt.Errorf("internal server error when saving world: %s", err)
	}

	return &rove.RegisterResponse{}, nil
}

// Rover returns rover information for a gRPC request
func (s *Server) Rover(ctx context.Context, req *rove.RoverRequest) (*rove.RoverResponse, error) {
	response := &rove.RoverResponse{}
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account name")

	} else if resp, err := s.accountant.GetValue(req.Account, "rover"); err != nil {
		return nil, err

	} else if rover, err := s.world.GetRover(resp); err != nil {
		return nil, fmt.Errorf("error getting rover: %s", err)

	} else {
		response = &rove.RoverResponse{
			Name: rover.Name,
			Position: &rove.Vector{
				X: int32(rover.Pos.X),
				Y: int32(rover.Pos.Y),
			},
			Range:     int32(rover.Range),
			Inventory: rover.Inventory,
			Integrity: int32(rover.Integrity),
		}
	}
	return response, nil
}

// Radar returns the radar information for a rover
func (s *Server) Radar(ctx context.Context, req *rove.RadarRequest) (*rove.RadarResponse, error) {
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account name")
	}

	response := &rove.RadarResponse{}

	resp, err := s.accountant.GetValue(req.Account, "rover")
	if err != nil {
		return nil, err

	} else if rover, err := s.world.GetRover(resp); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else if radar, err := s.world.RadarFromRover(resp); err != nil {
		return nil, fmt.Errorf("error getting radar from rover: %s", err)

	} else {
		response.Tiles = radar
		response.Range = int32(rover.Range)
	}

	return response, nil
}

// Commands issues commands to the world based on a gRPC request
func (s *Server) Commands(ctx context.Context, req *rove.CommandsRequest) (*rove.CommandsResponse, error) {
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account")
	}
	resp, err := s.accountant.GetValue(req.Account, "rover")
	if err != nil {
		return nil, err
	}

	var cmds []game.Command
	for _, c := range req.Commands {
		cmds = append(cmds, game.Command{
			Bearing: c.Bearing,
			Command: c.Command})
	}

	if err := s.world.Enqueue(resp, cmds...); err != nil {
		return nil, err
	}

	return &rove.CommandsResponse{}, nil
}
