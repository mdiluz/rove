package persistence

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Dummy struct {
	Success bool
	Value   int
}

func TestPersistence_LoadSave(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "rove_persistence_test")
	assert.NoError(t, err, "Failed to get tempdir path")

	assert.NoError(t, SetPath(tmp), "Failed to get set tempdir to persistence path")

	// Try and save out the dummy
	var dummy Dummy
	dummy.Success = true
	assert.NoError(t, Save("test", dummy), "Failed to save out dummy file")

	// Load back the dummy
	dummy = Dummy{}
	assert.NoError(t, Load("test", &dummy), "Failed to load in dummy file")
	assert.Equal(t, true, dummy.Success, "Did not successfully load true value from file")
}

func TestPersistence_LoadSaveAll(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "rove_persistence_test")
	assert.NoError(t, err, "Failed to get tempdir path")

	assert.NoError(t, SetPath(tmp), "Failed to get set tempdir to persistence path")

	// Try and save out the dummy
	var dummyA Dummy
	var dummyB Dummy
	dummyA.Value = 1
	dummyB.Value = 2
	assert.NoError(t, SaveAll("a", dummyA, "b", dummyB), "Failed to save out dummy file")

	// Load back the dummy
	dummyA = Dummy{}
	dummyB = Dummy{}
	assert.NoError(t, LoadAll("a", &dummyA, "b", &dummyB), "Failed to load in dummy file")
	assert.Equal(t, 1, dummyA.Value, "Did not successfully load int value from file")
	assert.Equal(t, 2, dummyB.Value, "Did not successfully load int value from file")
}
