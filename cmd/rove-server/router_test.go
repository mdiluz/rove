package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleStatus(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()

	HandleStatus(response, request)

	var status RouterStatus
	json.NewDecoder(response.Body).Decode(&status)

	if status.Ready != true {
		t.Errorf("got false for /status")
	}
}