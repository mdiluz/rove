// +build integration

package main

import (
	"testing"

	"github.com/mdiluz/rove/pkg/rove"
)

var serverUrl = "http://localhost:8080"

func TestServerStatus(t *testing.T) {
	conn := rove.NewConnection(serverUrl)
	status := conn.Status()
	if !status.Ready {
		t.Error("Server did not return that it was ready")
	}
}
