package rove

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Server is a simple wrapper to a server path
type Server string

// Get performs a Get request
func (s Server) Get(path string, out interface{}) error {
	u := url.URL{
		Scheme: "http",
		Host:   string(s),
		Path:   path,
	}
	if resp, err := http.Get(u.String()); err != nil {
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

// Post performs a Post request
func (s Server) Post(path string, in, out interface{}) error {
	u := url.URL{
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
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(marshalled))
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
