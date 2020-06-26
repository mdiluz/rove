// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/stretchr/testify/assert"
)

// Server is a simple wrapper to a server path
type Server string

// Request performs a HTTP
func (s Server) Request(method, path string, in, out interface{}) error {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:8080", string(s)),
		Path:   path,
	}
	client := &http.Client{}

	// Marshal the input
	marshalled, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// Set up the request
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	// Do the POST
	req.Header.Set("Content-Type", "application/json")
	if resp, err := client.Do(req); err != nil {
		return err

	} else if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body to code %d", resp.StatusCode)
		}
		return fmt.Errorf("http returned status %d: %s", resp.StatusCode, string(body))

	} else {
		return json.NewDecoder(resp.Body).Decode(out)
	}
}

var serv = Server(os.Getenv("ROVE_HTTP"))

func TestServer_Status(t *testing.T) {
	req := &rove.StatusRequest{}
	resp := &rove.StatusResponse{}
	if err := serv.Request("GET", "status", req, resp); err != nil {
		log.Fatal(err)
	}
}

func TestServer_Register(t *testing.T) {
	req := &rove.RegisterRequest{Name: uuid.New().String()}
	resp := &rove.RegisterResponse{}
	err := serv.Request("POST", "register", req, resp)
	assert.NoError(t, err, "First register attempt should pass")
	err = serv.Request("POST", "register", req, resp)
	assert.Error(t, err, "Second identical register attempt should fail")
}

func TestServer_Command(t *testing.T) {
	acc := uuid.New().String()
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &rove.RegisterResponse{})
	assert.NoError(t, err, "First register attempt should pass")

	err = serv.Request("POST", "commands", &rove.CommandsRequest{
		Account: acc,
		Commands: []*rove.Command{
			{
				Command: "move",
				Bearing: "NE",
			},
		},
	}, &rove.CommandsResponse{})
	assert.NoError(t, err, "Commands should should pass")
}

func TestServer_Radar(t *testing.T) {
	acc := uuid.New().String()
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &rove.RegisterResponse{})
	assert.NoError(t, err, "First register attempt should pass")

	resp := &rove.RadarResponse{}
	err = serv.Request("POST", "radar", &rove.RadarRequest{
		Account: acc,
	}, resp)
	assert.NoError(t, err, "Radar sould pass should pass")
	assert.NotZero(t, resp.Range, "Radar should return valid range")
	w := int(resp.Range*2 + 1)
	assert.Equal(t, w*w, len(resp.Tiles), "radar should return correct number of tiles")
}

func TestServer_Rover(t *testing.T) {
	acc := uuid.New().String()
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &rove.RegisterResponse{})
	assert.NoError(t, err, "First register attempt should pass")

	resp := &rove.RoverResponse{}
	err = serv.Request("POST", "rover", &rove.RoverRequest{
		Account: acc,
	}, resp)
	assert.NoError(t, err, "Rover sould pass should pass")
	assert.NotZero(t, resp.Range, "Rover should return valid range")
	assert.NotZero(t, len(resp.Name), "Rover should return valid name")
	assert.NotZero(t, resp.Position, "Rover should return valid position")
}
