// +build integration

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

var serverUrl = "localhost:80"

// Connection is the container for a simple connection to the server
type Connection struct {
	host string
}

// NewConnection sets up a new connection to a server host
func NewConnection(host string) *Connection {
	return &Connection{
		host: host,
	}
}

// Status returns the current status of the server
func (c *Connection) Status() (status StatusResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "status",
	}

	if resp, err := http.Get(url.String()); err != nil {
		return StatusResponse{}, err
	} else if resp.StatusCode != http.StatusOK {
		return StatusResponse{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&status)
	}

	return
}

// Register registers a new account on the server
func (c *Connection) Register(name string) (register RegisterResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "register",
	}

	// Marshal the register data struct
	data := RegisterData{Name: name}
	marshalled, err := json.Marshal(data)

	// Set up the request
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		return RegisterResponse{}, err
	} else {
		defer resp.Body.Close()

		// Handle any errors
		if resp.StatusCode != http.StatusOK {
			return RegisterResponse{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
		} else {
			// Decode the reply
			err = json.NewDecoder(resp.Body).Decode(&register)
		}
	}

	return
}

func TestStatus(t *testing.T) {
	conn := NewConnection(serverUrl)

	if status, err := conn.Status(); err != nil {
		t.Errorf("Status returned error: %s", err)
	} else if !status.Ready {
		t.Error("Server did not return that it was ready")
	} else if len(status.Version) == 0 {
		t.Error("Server returned blank version")
	}
}

func TestRegister(t *testing.T) {
	conn := NewConnection(serverUrl)

	a := uuid.New().String()
	reg1, err := conn.Register(a)
	if err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if !reg1.Success {
		t.Error("Server did not success for Register")
	} else if len(reg1.Id) == 0 {
		t.Error("Server returned empty registration ID")
	}

	b := uuid.New().String()
	reg2, err := conn.Register(b)
	if err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if !reg2.Success {
		t.Error("Server did not success for Register")
	} else if len(reg2.Id) == 0 {
		t.Error("Server returned empty registration ID")
	}

	if reg2, err := conn.Register(a); err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if reg2.Success {
		t.Error("Server should have failed to register duplicate name")
	}
}
