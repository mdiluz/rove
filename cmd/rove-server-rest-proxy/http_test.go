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

func TestServer_ServerStatus(t *testing.T) {
	req := &rove.ServerStatusRequest{}
	resp := &rove.ServerStatusResponse{}
	if err := serv.Request("GET", "server-status", req, resp); err != nil {
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
	var resp rove.RegisterResponse
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &resp)
	assert.NoError(t, err, "First register attempt should pass")

	req := &rove.CommandRequest{
		Account: &rove.Account{
			Name: resp.Account.Name,
		},
		Commands: []*rove.Command{
			{
				Command: "move",
				Bearing: "NE",
			},
		},
	}

	assert.Error(t, serv.Request("POST", "command", req, &rove.CommandResponse{}), "Commands should fail with no secret")

	req.Account.Secret = resp.Account.Secret
	assert.NoError(t, serv.Request("POST", "command", req, &rove.CommandResponse{}), "Commands should pass")
}

func TestServer_Radar(t *testing.T) {
	acc := uuid.New().String()
	var reg rove.RegisterResponse
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &reg)
	assert.NoError(t, err, "First register attempt should pass")

	resp := &rove.RadarResponse{}
	req := &rove.RadarRequest{
		Account: &rove.Account{
			Name: reg.Account.Name,
		},
	}

	assert.Error(t, serv.Request("POST", "radar", req, resp), "Radar should fail without secret")
	req.Account.Secret = reg.Account.Secret

	assert.NoError(t, serv.Request("POST", "radar", req, resp), "Radar should pass")
	assert.NotZero(t, resp.Range, "Radar should return valid range")

	w := int(resp.Range*2 + 1)
	assert.Equal(t, w*w, len(resp.Tiles), "radar should return correct number of tiles")
	assert.Equal(t, w*w, len(resp.Objects), "radar should return correct number of objects")
}

func TestServer_Status(t *testing.T) {
	acc := uuid.New().String()
	var reg rove.RegisterResponse
	err := serv.Request("POST", "register", &rove.RegisterRequest{Name: acc}, &reg)
	assert.NoError(t, err, "First register attempt should pass")

	resp := &rove.StatusResponse{}
	req := &rove.StatusRequest{
		Account: &rove.Account{
			Name: reg.Account.Name,
		},
	}

	assert.Error(t, serv.Request("POST", "status", req, resp), "Status should fail without secret")
	req.Account.Secret = reg.Account.Secret

	assert.NoError(t, serv.Request("POST", "status", req, resp), "Status should pass")
	assert.NotZero(t, resp.Range, "Rover should return valid range")
	assert.NotZero(t, len(resp.Name), "Rover should return valid name")
	assert.NotZero(t, resp.Position, "Rover should return valid position")
	assert.NotZero(t, resp.Integrity, "Rover should have positive integrity")
}
