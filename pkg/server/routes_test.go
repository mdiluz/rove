package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/stretchr/testify/assert"
)

func TestHandleStatus(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()

	s := NewServer()
	s.Initialise()
	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

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
	s.Initialise()
	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RegisterResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /register")
	}
}

func TestHandleSpawn(t *testing.T) {
	s := NewServer()
	s.Initialise()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")
	data := rove.SpawnData{}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, path.Join("/", a.Id.String(), "/spawn"), bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.SpawnResponse
	json.NewDecoder(response.Body).Decode(&status)
	assert.Equal(t, http.StatusOK, response.Code)

	if status.Success != true {
		t.Errorf("got false for /spawn: %s", status.Error)
	}
}

func TestHandleCommand(t *testing.T) {
	s := NewServer()
	s.Initialise()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, inst, err := s.SpawnRoverForAccount(a.Id)

	pos, err := s.world.RoverPosition(inst)
	assert.NoError(t, err, "Couldn't get rover position")

	data := rove.CommandData{
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

	request, _ := http.NewRequest(http.MethodPost, path.Join("/", a.Id.String(), "/command"), bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.CommandResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /command: %s", status.Error)
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
	s.Initialise()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, _, err = s.SpawnRoverForAccount(a.Id)

	request, _ := http.NewRequest(http.MethodGet, path.Join("/", a.Id.String(), "/radar"), nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RadarResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /radar: %s", status.Error)
	}

	// TODO: Verify the radar information
}

func TestHandleRover(t *testing.T) {
	s := NewServer()
	s.Initialise()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, _, err = s.SpawnRoverForAccount(a.Id)

	request, _ := http.NewRequest(http.MethodGet, path.Join("/", a.Id.String(), "/rover"), nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RoverResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /rover: %s", status.Error)
	}

	// TODO: Verify the radar information
}
