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

	"github.com/mdiluz/rove/pkg/rove"
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
}

func TestServer_Command(t *testing.T) {
}

func TestServer_Radar(t *testing.T) {
}

func TestServer_Rover(t *testing.T) {
}
