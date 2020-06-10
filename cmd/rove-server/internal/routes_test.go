package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/vector"
	"github.com/stretchr/testify/assert"
)

func TestHandleStatus(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()

	s := NewServer()
	s.Initialise(true)
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
	s.Initialise(true)
	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RegisterResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /register")
	}
}

func TestHandleCommand(t *testing.T) {
	s := NewServer()
	s.Initialise(false) // Leave the world empty with no obstacles
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	_, inst, err := s.SpawnRoverForAccount(a.Id)
	assert.NoError(t, s.world.WarpRover(inst, vector.Vector{}))

	attribs, err := s.world.RoverAttributes(inst)
	assert.NoError(t, err, "Couldn't get rover position")

	data := rove.CommandData{
		Commands: []game.Command{
			{
				Command:  game.CommandMove,
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

	// Tick the command queues to progress the move command
	s.world.EnqueueAllIncoming()
	s.world.ExecuteCommandQueues()

	attribs2, err := s.world.RoverAttributes(inst)
	assert.NoError(t, err, "Couldn't get rover position")
	attribs.Pos.Add(vector.Vector{X: 0.0, Y: attrib.Speed * 1}) // Should have moved north by the speed and duration
	assert.Equal(t, attribs.Pos, attribs2.Pos, "Rover should have moved by bearing")
}

func TestHandleRadar(t *testing.T) {
	s := NewServer()
	s.Initialise(false) // Spawn a clean world
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn the rover rover for the account
	attrib, id, err := s.SpawnRoverForAccount(a.Id)
	assert.NoError(t, err)

	// Warp this rover to 0,0
	assert.NoError(t, s.world.WarpRover(id, vector.Vector{}))

	// Explicity set a few nearby tiles
	wallPos1 := vector.Vector{X: 0, Y: -1}
	wallPos2 := vector.Vector{X: 1, Y: 1}
	rockPos := vector.Vector{X: 1, Y: 3}
	emptyPos := vector.Vector{X: -2, Y: -3}
	assert.NoError(t, s.world.Atlas.SetTile(wallPos1, game.TileWall))
	assert.NoError(t, s.world.Atlas.SetTile(wallPos2, game.TileWall))
	assert.NoError(t, s.world.Atlas.SetTile(rockPos, game.TileRock))
	assert.NoError(t, s.world.Atlas.SetTile(emptyPos, game.TileEmpty))

	request, _ := http.NewRequest(http.MethodGet, path.Join("/", a.Id.String(), "/radar"), nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RadarResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /radar: %s", status.Error)
	}

	scope := attrib.Range*2 + 1
	radarOrigin := vector.Vector{X: -attrib.Range, Y: -attrib.Range}

	// Make sure the rover tile is correct
	assert.Equal(t, game.TileRover, status.Tiles[len(status.Tiles)/2])

	// Check our other tiles
	wallPos1.Add(radarOrigin.Negated())
	wallPos2.Add(radarOrigin.Negated())
	rockPos.Add(radarOrigin.Negated())
	emptyPos.Add(radarOrigin.Negated())
	assert.Equal(t, game.TileWall, status.Tiles[wallPos1.X+wallPos1.Y*scope])
	assert.Equal(t, game.TileWall, status.Tiles[wallPos2.X+wallPos2.Y*scope])
	assert.Equal(t, game.TileRock, status.Tiles[rockPos.X+rockPos.Y*scope])
	assert.Equal(t, game.TileEmpty, status.Tiles[emptyPos.X+emptyPos.Y*scope])

}

func TestHandleRover(t *testing.T) {
	s := NewServer()
	s.Initialise(true)
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")

	// Spawn one rover for the account
	attribs, _, err := s.SpawnRoverForAccount(a.Id)
	assert.NoError(t, err)

	request, _ := http.NewRequest(http.MethodGet, path.Join("/", a.Id.String(), "/rover"), nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var status rove.RoverResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /rover: %s", status.Error)
	} else if attribs != status.Attributes {
		t.Errorf("Missmatched attributes: %+v, !=%+v", attribs, status.Attributes)
	}
}
