package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
	"google.golang.org/grpc"
)

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

func (s *Server) Register(ctx context.Context, req *rove.RegisterRequest) (*rove.RegisterResponse, error) {
	if len(req.Name) == 0 {
		return nil, fmt.Errorf("empty account name")
	}

	if _, err := s.accountant.Register(ctx, &accounts.RegisterInfo{Name: req.Name}, grpc.WaitForReady(true)); err != nil {
		return nil, err

	} else if _, _, err := s.SpawnRoverForAccount(req.Name); err != nil {
		return nil, fmt.Errorf("failed to spawn rover for account: %s", err)

	} else if err := s.SaveWorld(); err != nil {
		return nil, fmt.Errorf("internal server error when saving world: %s", err)

	}

	return &rove.RegisterResponse{}, nil
}

func (s *Server) Rover(ctx context.Context, req *rove.RoverRequest) (*rove.RoverResponse, error) {
	response := &rove.RoverResponse{}
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account name")

	} else if resp, err := s.accountant.GetValue(ctx, &accounts.DataKey{Account: req.Account, Key: "rover"}); err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)

	} else if id, err := uuid.Parse(resp.Value); err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)

	} else if attrib, err := s.world.RoverAttributes(id); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else if pos, err := s.world.RoverPosition(id); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else {
		response = &rove.RoverResponse{
			Name: attrib.Name,
			Position: &rove.Vector{
				X: int32(pos.X),
				Y: int32(pos.Y),
			},
			Range:    int32(attrib.Range),
			Capacity: int32(attrib.Capacity),
		}
	}
	return response, nil
}

func (s *Server) Radar(ctx context.Context, req *rove.RadarRequest) (*rove.RadarResponse, error) {
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account name")
	}

	response := &rove.RadarResponse{}

	resp, err := s.accountant.GetValue(ctx, &accounts.DataKey{Account: req.Account, Key: "rover"})
	if err != nil {
		return nil, fmt.Errorf("gRPC failed to contact accountant: %s", err)
	}

	if id, err := uuid.Parse(resp.Value); err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)

	} else if attrib, err := s.world.RoverAttributes(id); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else if radar, err := s.world.RadarFromRover(id); err != nil {
		return nil, fmt.Errorf("error getting radar from rover: %s", err)

	} else {
		response.Tiles = radar
		response.Range = int32(attrib.Range)
	}

	return response, nil
}

func (s *Server) Commands(ctx context.Context, req *rove.CommandsRequest) (*rove.CommandsResponse, error) {
	if len(req.Account) == 0 {
		return nil, fmt.Errorf("empty account")
	}
	resp, err := s.accountant.GetValue(ctx, &accounts.DataKey{Account: req.Account, Key: "rover"})

	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(resp.Value)
	if err != nil {
		return nil, fmt.Errorf("account had invalid rover ID: %s", resp.Value)
	}

	var cmds []game.Command
	for _, c := range req.Commands {
		cmds = append(cmds, game.Command{
			Bearing: c.Bearing,
			Command: c.Command})
	}

	if err := s.world.Enqueue(id, cmds...); err != nil {
		return nil, err
	}

	return &rove.CommandsResponse{}, nil
}
