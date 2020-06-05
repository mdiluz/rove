package main

import (
	"flag"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	s := server.NewServer(server.OptionPort(8080))
	s.Initialise()
	go s.Run()

	code := m.Run()

	s.Close()

	os.Exit(code)
}

func Test_InnerMain(t *testing.T) {
	// Set up the flags to act locally and use a temporary file
	flag.Set("data", path.Join(os.TempDir(), uuid.New().String()))

	// First attempt should error
	assert.Error(t, InnerMain("status"))

	// Now set the host
	flag.Set("host", "localhost:8080")

	// No error now as we have a host
	assert.NoError(t, InnerMain("status"))

	// Register should fail without a name
	assert.Error(t, InnerMain("register"))

	// These methods should fail without an account
	assert.Error(t, InnerMain("spawn"))
	assert.Error(t, InnerMain("command"))
	assert.Error(t, InnerMain("radar"))
	assert.Error(t, InnerMain("rover"))

	// Now set the name
	flag.Set("name", uuid.New().String())

	// Perform the register
	assert.NoError(t, InnerMain("register"))

	// We've not spawned a rover yet so these should fail
	assert.Error(t, InnerMain("command")) // Currently not erroring, needs investigation
	assert.Error(t, InnerMain("radar"))
	assert.Error(t, InnerMain("rover"))

	// Spawn a rover
	assert.NoError(t, InnerMain("spawn"))

	// These should now work
	assert.NoError(t, InnerMain("command"))
	assert.NoError(t, InnerMain("radar"))
	assert.NoError(t, InnerMain("rover"))
}
