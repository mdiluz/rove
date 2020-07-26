// +build integration

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_InnerMain(t *testing.T) {
	os.Setenv("NO_TLS", "1")

	// Use temporary local user data
	tmp, err := ioutil.TempDir(os.TempDir(), "rove-")
	assert.NoError(t, err)
	os.Setenv("ROVE_USER_DATA", tmp)

	// Used for configuring this test
	var address = os.Getenv("ROVE_GRPC")
	if len(address) == 0 {
		log.Fatal("Must set $ROVE_GRPC")
	}

	// First attempt should error without a host
	assert.Error(t, InnerMain("server-status"))

	// Set the host in the config
	assert.NoError(t, InnerMain("config", address))
	assert.NoError(t, InnerMain("server-status"))

	// Register should fail without a name
	assert.Error(t, InnerMain("register"))

	// These methods should fail without an account
	assert.Error(t, InnerMain("radar"))
	assert.Error(t, InnerMain("status"))

	// Now set the name
	assert.NoError(t, InnerMain("register", uuid.New().String()))

	// These should now work
	assert.NoError(t, InnerMain("radar"))
	assert.NoError(t, InnerMain("status"))

	// Commands should fail with no commands
	assert.Error(t, InnerMain("command"))

	// Give it commands
	assert.NoError(t, InnerMain("command", "toggle"))
	assert.NoError(t, InnerMain("command", "stash"))
	assert.NoError(t, InnerMain("command", "repair"))
	assert.NoError(t, InnerMain("command", "broadcast", "abc"))
	assert.NoError(t, InnerMain("command", "10", "wait"))
	assert.NoError(t, InnerMain("command", "1", "wait", "turn", "NW", "toggle", "broadcast", "zyx"))

	// Give it malformed commands
	assert.Error(t, InnerMain("command", "unknown"))
	assert.Error(t, InnerMain("command", "broadcast"))
	assert.Error(t, InnerMain("command", "1"))
}
