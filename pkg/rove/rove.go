package rove

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Connection struct {
	host string
}

func NewConnection(host string) *Connection {
	return &Connection{
		host: host,
	}
}

// ServerStatus is a struct that contains information on the status of the server
type ServerStatus struct {
	Ready bool `json:"ready"`
}

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
