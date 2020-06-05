package rove

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Server is a simple wrapper to a server path
type Server string

// GET performs a GET request
func (s Server) GET(path string, out interface{}) error {
	url := url.URL{
		Scheme: "http",
		Host:   string(s),
		Path:   path,
	}
	if resp, err := http.Get(url.String()); err != nil {
		return err

	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Get returned status %d: %s", resp.StatusCode, resp.Status)

	} else {
		return json.NewDecoder(resp.Body).Decode(out)
	}
}

// POST performs a POST request
func (s Server) POST(path string, in interface{}, out interface{}) error {
	url := url.URL{
		Scheme: "http",
		Host:   string(s),
		Path:   path,
	}
	client := &http.Client{}

	// Marshal the input
	marshalled, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// Set up the request
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	// Do the POST
	req.Header.Set("Content-Type", "application/json")
	if resp, err := client.Do(req); err != nil {
		return err

	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http returned status %d", resp.StatusCode)

	} else {
		return json.NewDecoder(resp.Body).Decode(out)
	}
}
