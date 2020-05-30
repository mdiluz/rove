// +build integration

package rove

import (
	"testing"
)

var serverUrl = "localhost:8080"

func TestStatus(t *testing.T) {
	conn := NewConnection(serverUrl)
	if status, err := conn.Status(); err != nil {
		t.Errorf("Status returned error: %s", err)
	} else if !status.Ready {
		t.Error("Server did not return that it was ready")
	}
}

func TestRegister(t *testing.T) {
	conn := NewConnection(serverUrl)
	if reg, err := conn.Register(); err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if !reg.Success {
		t.Error("Server did not success for Register")
	} else if len(reg.Id) == 0 {
		t.Error("Server returned empty registration ID")
	}
}
