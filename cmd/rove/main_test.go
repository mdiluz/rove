// +build integration

package main

import (
	"flag"
	"log"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_InnerMain(t *testing.T) {

	var address = os.Getenv("ROVE_SERVER_ADDRESS")
	if len(address) == 0 {
		log.Fatal("Must set ROVE_SERVER_ADDRESS")
	}

	// Set up the flags to act locally and use a temporary file
	flag.Set("data", path.Join(os.TempDir(), uuid.New().String()))

	// First attempt should error
	assert.Error(t, InnerMain("status"))

	// Now set the host
	flag.Set("host", address)

	// No error now as we have a host
	assert.NoError(t, InnerMain("status"))

	// Register should fail without a name
	assert.Error(t, InnerMain("register"))

	// These methods should fail without an account
	assert.Error(t, InnerMain("move"))
	assert.Error(t, InnerMain("radar"))
	assert.Error(t, InnerMain("rover"))

	// Now set the name
	flag.Set("name", uuid.New().String())

	// Perform the register
	assert.NoError(t, InnerMain("register"))

	// These should now work
	assert.NoError(t, InnerMain("radar"))
	assert.NoError(t, InnerMain("rover"))

	// Move should work with arguments
	flag.Set("bearing", "N")
	flag.Set("duration", "1")
	assert.NoError(t, InnerMain("move"))
}
