package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/stretchr/testify/assert"
)

func TestHandleStatus(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()

	s := NewServer()
	s.wrapHandler(http.MethodGet, HandleStatus)(response, request)

	var status rove.StatusResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Ready != true {
		t.Errorf("got false for /status")
	}

	if len(status.Version) == 0 {
		t.Errorf("got empty version info")
	}
}

func TestHandleRegister(t *testing.T) {
	data := rove.RegisterData{Name: "one"}
	b, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s := NewServer()
	s.wrapHandler(http.MethodPost, HandleRegister)(response, request)

	var status rove.RegisterResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /register")
	}
}

func TestHandleSpawn(t *testing.T) {
	s := NewServer()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")
	data := rove.SpawnData{Id: a.Id.String()}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, "/spawn", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.wrapHandler(http.MethodPost, HandleSpawn)(response, request)

	var status rove.SpawnResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /spawn")
	}
}

func TestHandleCommand(t *testing.T) {
	s := NewServer()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, inst, err := s.SpawnRoverForAccount(a.Id)

	pos, err := s.world.RoverPosition(inst)
	assert.NoError(t, err, "Couldn't get rover position")

	data := rove.CommandData{
		Id: a.Id.String(),
		Commands: []rove.Command{
			{
				Command:  rove.CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, "/command", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.wrapHandler(http.MethodPost, HandleCommand)(response, request)

	var status rove.CommandResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /command")
	}

	attrib, err := s.world.RoverAttributes(inst)
	assert.NoError(t, err, "Couldn't get rover attribs")

	pos2, err := s.world.RoverPosition(inst)
	assert.NoError(t, err, "Couldn't get rover position")
	pos.Add(game.Vector{X: 0.0, Y: attrib.Speed * 1}) // Should have moved north by the speed and duration
	assert.Equal(t, pos, pos2, "Rover should have moved by bearing")
}

func TestHandleRadar(t *testing.T) {
	s := NewServer()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, _, err = s.SpawnRoverForAccount(a.Id)

	data := rove.RadarData{
		Id: a.Id.String(),
	}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, "/radar", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.wrapHandler(http.MethodPost, HandleRadar)(response, request)

	var status rove.RadarResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /radar")
	}

	// TODO: Verify the radar information
}

func TestHandleRover(t *testing.T) {
	s := NewServer()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, _, err = s.SpawnRoverForAccount(a.Id)

	data := rove.RoverData{
		Id: a.Id.String(),
	}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, "/rover", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.wrapHandler(http.MethodPost, HandleRover)(response, request)

	var status rove.RoverResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /rover")
	}

	// TODO: Verify the radar information
}
