// +build integration

package main

import (
	"testing"

	"github.com/mdiluz/rove/pkg/rove"
)

var serverUrl = "localhost:8080"

func TestServerStatus(t *testing.T) {
	conn := rove.NewConnection(serverUrl)
	if status, err := conn.Status(); err != nil {
		t.Errorf("Status returned error: %s", err)
	} else if !status.Ready {
		t.Error("Server did not return that it was ready")
	}
}
