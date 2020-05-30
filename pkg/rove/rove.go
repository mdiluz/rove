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

// ServerStatus is a struct that contains information on the status of the server
type ServerStatus struct {
	Ready bool `json:"ready"`
}

// Status returns the current status of the server
func (c *Connection) Status() (status ServerStatus, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "status",
	}

	if resp, err := http.Get(url.String()); err != nil {
		return ServerStatus{}, err
	} else if resp.StatusCode != http.StatusOK {
		return ServerStatus{}, fmt.Errorf("Status request returned %d", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&status)
	}

	return
}
