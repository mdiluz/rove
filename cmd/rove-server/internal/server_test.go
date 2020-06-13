package internal

import (
	"os"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Error("Failed to create server")
	}
}

func TestNewServer_OptionAddress(t *testing.T) {
	server := NewServer(OptionAddress(":1234"))
	if server == nil {
		t.Error("Failed to create server")
	} else if server.address != ":1234" {
		t.Error("Failed to set server address")
	}
}

func TestNewServer_OptionPersistentData(t *testing.T) {
	server := NewServer(OptionPersistentData())
	if server == nil {
		t.Error("Failed to create server")
	} else if server.persistence != PersistentData {
		t.Error("Failed to set server persistent data")
	}
}

func TestServer_Run(t *testing.T) {
	os.Setenv("ROVE_ACCOUNTANT_GRPC", "n/a")
	server := NewServer()
	if server == nil {
		t.Error("Failed to create server")
	} else if err := server.Initialise(true); err != nil {
		t.Error(err)
	}

	go server.Run()

	if err := server.StopAndClose(); err != nil {
		t.Error(err)
	}
}

func TestServer_RunPersistentData(t *testing.T) {
	os.Setenv("ROVE_ACCOUNTANT_GRPC", "n/a")
	server := NewServer(OptionPersistentData())
	if server == nil {
		t.Error("Failed to create server")
	} else if err := server.Initialise(true); err != nil {
		t.Error(err)
	}

	go server.Run()

	if err := server.StopAndClose(); err != nil {
		t.Error(err)
	}
}
