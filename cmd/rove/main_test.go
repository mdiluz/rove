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
	assert.Error(t, InnerMain("status"))

	// Set the host in the config
	assert.NoError(t, InnerMain("config", address))
	assert.NoError(t, InnerMain("status"))

	// Register should fail without a name
	assert.Error(t, InnerMain("register"))

	// These methods should fail without an account
	assert.Error(t, InnerMain("radar"))
	assert.Error(t, InnerMain("rover"))

	// Now set the name
	assert.NoError(t, InnerMain("register", uuid.New().String()))

	// These should now work
	assert.NoError(t, InnerMain("radar"))
	assert.NoError(t, InnerMain("rover"))

	// Commands should fail with no commands
	assert.Error(t, InnerMain("commands"))

	// Give it commands
	assert.NoError(t, InnerMain("commands", "move", "N"))
	assert.NoError(t, InnerMain("commands", "stash"))
	assert.NoError(t, InnerMain("commands", "repair"))

	// Give it malformed commands
	assert.Error(t, InnerMain("commands", "move", "stash"))
}
