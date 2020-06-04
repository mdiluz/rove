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
	"github.com/stretchr/testify/assert"
)

var serverUrl = "localhost:80"

func TestStatus(t *testing.T) {
	url := url.URL{
		Scheme: "http",
		Host:   serverUrl,
		Path:   "status",
	}
	resp, err := http.Get(url.String())
	assert.NoError(t, err, "http.Get must not return error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "http.Get must return StatusOK")

	var status StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&status)
	assert.NoError(t, err, "json decode must not return error")

	assert.NoError(t, err, "Status must not return error")
	assert.True(t, status.Ready, "Server must return ready")
	assert.NotZero(t, len(status.Version), "Version must not be empty")
}

// helper for register test
func register(name string) (register RegisterResponse, err error) {
	url := url.URL{
		Scheme: "http",
		Host:   serverUrl,
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

func TestRegister(t *testing.T) {
	a := uuid.New().String()
	reg1, err := register(a)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, reg1.Success, "Register must return success")
	assert.NotZero(t, len(reg1.Id), "Register must return registration ID")

	b := uuid.New().String()
	reg2, err := register(b)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, reg2.Success, "Register must return success")
	assert.NotZero(t, len(reg2.Id), "Register must return registration ID")

	reg2, err = register(a)
	assert.NoError(t, err, "Register must not return error")
	assert.False(t, reg2.Success, "Register must return fail for duplicate registration")
}
