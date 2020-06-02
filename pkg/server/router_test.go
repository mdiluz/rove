package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleStatus(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()

	s := NewServer()
	s.HandleStatus(response, request)

	var status StatusResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Ready != true {
		t.Errorf("got false for /status")
	}

	if len(status.Version) == 0 {
		t.Errorf("got empty version info")
	}
}

func TestHandleRegister(t *testing.T) {
	data := RegisterData{Name: "one"}
	b, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s := NewServer()
	s.HandleRegister(response, request)

	var status RegisterResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /register")
	}
}

func TestHandleSpawn(t *testing.T) {
	s := NewServer()
	a, err := s.accountant.RegisterAccount("test")
	assert.NoError(t, err, "Error registering account")
	data := SpawnData{BasicAccountData: BasicAccountData{Id: a.Id.String()}}

	b, err := json.Marshal(data)
	assert.NoError(t, err, "Error marshalling data")

	request, _ := http.NewRequest(http.MethodPost, "/spawn", bytes.NewReader(b))
	response := httptest.NewRecorder()

	s.HandleSpawn(response, request)

	var status SpawnResponse
	json.NewDecoder(response.Body).Decode(&status)

	if status.Success != true {
		t.Errorf("got false for /spawn")
	}
}
