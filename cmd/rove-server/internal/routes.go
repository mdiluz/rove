package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/version"
)

// ServerStatus returns the status of the current server to a gRPC request
func (s *Server) ServerStatus(context.Context, *rove.ServerStatusRequest) (*rove.ServerStatusResponse, error) {
	response := &rove.ServerStatusResponse{
		Ready:       true,
		Version:     version.Version,
		TickRate:    int32(s.minutesPerTick),
		CurrentTick: int32(s.world.CurrentTicks),
	}

	// If there's a schedule, respond with it
	if len(s.schedule.Entries()) > 0 {
		response.NextTick = s.schedule.Entries()[0].Next.Format("15:04:05")
	}

	return response, nil
}

// Register registers a new account for a gRPC request
func (s *Server) Register(ctx context.Context, req *rove.RegisterRequest) (*rove.RegisterResponse, error) {
	log.Printf("Handling register request: %s\n", req.Name)

	if len(req.Name) == 0 {
		return nil, fmt.Errorf("empty account name")
	}

	if acc, err := s.accountant.RegisterAccount(req.Name); err != nil {
		return nil, err

	} else if _, err := s.SpawnRoverForAccount(req.Name); err != nil {
		return nil, fmt.Errorf("failed to spawn rover for account: %s", err)

	} else if err := s.SaveWorld(); err != nil {
		return nil, fmt.Errorf("internal server error when saving world: %s", err)

	} else {
		return &rove.RegisterResponse{
			Account: &rove.Account{
				Name:   acc.Name,
				Secret: acc.Data["secret"],
			},
		}, nil
	}
}

// Status returns rover information for a gRPC request
func (s *Server) Status(ctx context.Context, req *rove.StatusRequest) (response *rove.StatusResponse, err error) {
	log.Printf("Handling status request: %s\n", req.Account.Name)

	if valid, err := s.accountant.VerifySecret(req.Account.Name, req.Account.Secret); err != nil {
		return nil, err

	} else if !valid {
		return nil, fmt.Errorf("Secret incorrect for account %s", req.Account.Name)

	} else if resp, err := s.accountant.GetValue(req.Account.Name, "rover"); err != nil {
		return nil, err

	} else if rover, err := s.world.GetRover(resp); err != nil {
		return nil, fmt.Errorf("error getting rover: %s", err)

	} else {
		var inv []byte
		for _, i := range rover.Inventory {
			inv = append(inv, byte(i.Type))
		}

		i, q := s.world.RoverCommands(resp)
		var incoming, queued []*rove.Command
		for _, i := range i {
			incoming = append(incoming, &rove.Command{
				Command: i.Command,
				Bearing: i.Bearing,
				Message: i.Message,
			})
		}
		for _, q := range q {
			queued = append(queued, &rove.Command{
				Command: q.Command,
				Bearing: q.Bearing,
				Message: q.Message,
			})
		}
		var logs []*rove.Log
		for _, log := range rover.Logs {
			logs = append(logs, &rove.Log{
				Text: log.Text,
				Time: fmt.Sprintf("%d", log.Time.Unix()), // proto uses strings under the hood for 64bit ints anyway
			})
		}

		response = &rove.StatusResponse{
			Name: rover.Name,
			Position: &rove.Vector{
				X: int32(rover.Pos.X),
				Y: int32(rover.Pos.Y),
			},
			Range:            int32(rover.Range),
			Inventory:        inv,
			Capacity:         int32(rover.Capacity),
			Integrity:        int32(rover.Integrity),
			MaximumIntegrity: int32(rover.MaximumIntegrity),
			Charge:           int32(rover.Charge),
			MaximumCharge:    int32(rover.MaximumCharge),
			IncomingCommands: incoming,
			QueuedCommands:   queued,
			Logs:             logs,
		}
	}
	return response, nil
}

// Radar returns the radar information for a rover
func (s *Server) Radar(ctx context.Context, req *rove.RadarRequest) (*rove.RadarResponse, error) {
	log.Printf("Handling radar request: %s\n", req.Account.Name)

	if valid, err := s.accountant.VerifySecret(req.Account.Name, req.Account.Secret); err != nil {
		return nil, err

	} else if !valid {
		return nil, fmt.Errorf("Secret incorrect for account %s", req.Account.Name)
	}

	response := &rove.RadarResponse{}

	resp, err := s.accountant.GetValue(req.Account.Name, "rover")
	if err != nil {
		return nil, err

	} else if rover, err := s.world.GetRover(resp); err != nil {
		return nil, fmt.Errorf("error getting rover attributes: %s", err)

	} else if radar, objs, err := s.world.RadarFromRover(resp); err != nil {
		return nil, fmt.Errorf("error getting radar from rover: %s", err)

	} else {
		response.Objects = objs
		response.Tiles = radar
		response.Range = int32(rover.Range)
	}

	return response, nil
}

// Command issues commands to the world based on a gRPC request
func (s *Server) Command(ctx context.Context, req *rove.CommandRequest) (*rove.CommandResponse, error) {
	log.Printf("Handling command request: %s and %+v\n", req.Account.Name, req.Commands)

	if valid, err := s.accountant.VerifySecret(req.Account.Name, req.Account.Secret); err != nil {
		return nil, err

	} else if !valid {
		return nil, fmt.Errorf("Secret incorrect for account %s", req.Account.Name)
	}

	resp, err := s.accountant.GetValue(req.Account.Name, "rover")
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

	return &rove.CommandResponse{}, nil
}
