package rove

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var server Server = "localhost:80"

func TestServer_Status(t *testing.T) {
	status, err := server.Status()
	assert.NoError(t, err, "Status must not return error")
	assert.True(t, status.Ready, "Server must return ready")
	assert.NotZero(t, len(status.Version), "Version must not be empty")
}

func TestServer_Register(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r1.Success, "Register must return success")
	assert.NotZero(t, len(r1.Id), "Register must return registration ID")

	d2 := RegisterData{
		Name: uuid.New().String(),
	}
	r2, err := server.Register(d2)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r2.Success, "Register must return success")
	assert.NotZero(t, len(r2.Id), "Register must return registration ID")

	r3, err := server.Register(d1)
	assert.NoError(t, err, "Register must not return error")
	assert.False(t, r3.Success, "Register must return fail for duplicate registration")
}

func TestServer_Spawn(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r1.Success, "Register must return success")
	assert.NotZero(t, len(r1.Id), "Register must return registration ID")

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r2.Success, "Register must return success")
}

func TestServer_Commands(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r1.Success, "Register must return success")
	assert.NotZero(t, len(r1.Id), "Register must return registration ID")

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r2.Success, "Register must return success")

	c := CommandsData{
		Id: r1.Id,
		Commands: []Command{
			{
				Command:  CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}
	r3, err := server.Commands(c)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r3.Success, "Register must return success")
}

func TestServer_Radar(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r1.Success, "Register must return success")
	assert.NotZero(t, len(r1.Id), "Register must return registration ID")

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r2.Success, "Register must return success")

	r := RadarData{
		Id: r1.Id,
	}
	r3, err := server.Radar(r)
	assert.NoError(t, err, "Register must not return error")
	assert.True(t, r3.Success, "Register must return success")
}
