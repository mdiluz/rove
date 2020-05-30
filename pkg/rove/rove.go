package rove

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

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

// StatusResponse is a struct that contains information on the status of the server
type StatusResponse struct {
	Ready bool `json:"ready"`
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

// RegisterResponse
type RegisterResponse struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
}

// Register registers a new player on the server
func (c *Connection) Register() (register RegisterResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "register",
	}

	if resp, err := http.Get(url.String()); err != nil {
		return RegisterResponse{}, err
	} else if resp.StatusCode != http.StatusOK {
		return RegisterResponse{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&register)
	}

	return
}
