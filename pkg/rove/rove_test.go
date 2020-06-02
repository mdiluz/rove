// +build integration

package rove

import (
	"testing"

	"github.com/google/uuid"
)

var serverUrl = "localhost:80"

func TestStatus(t *testing.T) {
	conn := NewConnection(serverUrl)

	if status, err := conn.Status(); err != nil {
		t.Errorf("Status returned error: %s", err)
	} else if !status.Ready {
		t.Error("Server did not return that it was ready")
	} else if len(status.Version) == 0 {
		t.Error("Server returned blank version")
	}
}

func TestRegister(t *testing.T) {
	conn := NewConnection(serverUrl)

	a := uuid.New().String()
	reg1, err := conn.Register(a)
	if err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if !reg1.Success {
		t.Error("Server did not success for Register")
	} else if len(reg1.Id) == 0 {
		t.Error("Server returned empty registration ID")
	}

	b := uuid.New().String()
	reg2, err := conn.Register(b)
	if err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if !reg2.Success {
		t.Error("Server did not success for Register")
	} else if len(reg2.Id) == 0 {
		t.Error("Server returned empty registration ID")
	}

	if reg2, err := conn.Register(a); err != nil {
		t.Errorf("Register returned error: %s", err)
	} else if reg2.Success {
		t.Error("Server should have failed to register duplicate name")
	}
}
