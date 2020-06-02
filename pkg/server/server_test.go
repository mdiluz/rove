package server

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Error("Failed to create server")
	}
}

func TestNewServer_OptionPort(t *testing.T) {
	server := NewServer(OptionPort(1234))
	if server == nil {
		t.Error("Failed to create server")
	} else if server.port != 1234 {
		t.Error("Failed to set server port")
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
	server := NewServer()
	if server == nil {
		t.Error("Failed to create server")
	} else if err := server.Initialise(); err != nil {
		t.Error(err)
	}
	go server.Run()

	if err := server.Close(); err != nil {
		t.Error(err)
	}
}

func TestServer_RunPersistentData(t *testing.T) {
	server := NewServer(OptionPersistentData())
	if server == nil {
		t.Error("Failed to create server")
	} else if err := server.Initialise(); err != nil {
		t.Error(err)
	}
	go server.Run()

	if err := server.Close(); err != nil {
		t.Error(err)
	}
}
