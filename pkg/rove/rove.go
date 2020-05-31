package rove

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mdiluz/rove/pkg/server"
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

// Status returns the current status of the server
func (c *Connection) Status() (status server.StatusResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "status",
	}

	if resp, err := http.Get(url.String()); err != nil {
		return server.StatusResponse{}, err
	} else if resp.StatusCode != http.StatusOK {
		return server.StatusResponse{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&status)
	}

	return
}

// Register registers a new account on the server
func (c *Connection) Register(name string) (register server.RegisterResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "register",
	}

	// Marshal the register data struct
	data := server.RegisterData{Name: name}
	marshalled, err := json.Marshal(data)

	// Set up the request
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		return server.RegisterResponse{}, err
	} else {
		defer resp.Body.Close()

		// Handle any errors
		if resp.StatusCode != http.StatusOK {
			return server.RegisterResponse{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
		} else {
			// Decode the reply
			err = json.NewDecoder(resp.Body).Decode(&register)
		}
	}

	return
}
