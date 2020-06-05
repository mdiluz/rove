package rove

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var server Server = "localhost:80"

func TestServer_Status(t *testing.T) {
	status, err := server.Status()
	assert.NoError(t, err)
	assert.True(t, status.Ready)
	assert.NotZero(t, len(status.Version))
}

func TestServer_Register(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	d2 := RegisterData{
		Name: uuid.New().String(),
	}
	r2, err := server.Register(d2)
	assert.NoError(t, err)
	assert.True(t, r2.Success)
	assert.NotZero(t, len(r2.Id))

	r3, err := server.Register(d1)
	assert.NoError(t, err)
	assert.False(t, r3.Success)
}

func TestServer_Spawn(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)
}

func TestServer_Command(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	c := CommandData{
		Id: r1.Id,
		Commands: []Command{
			{
				Command:  CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}
	r3, err := server.Command(c)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}

func TestServer_Radar(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	r := RadarData{
		Id: r1.Id,
	}
	r3, err := server.Radar(r)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}

func TestServer_Rover(t *testing.T) {
	d1 := RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := server.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := SpawnData{
		Id: r1.Id,
	}
	r2, err := server.Spawn(s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	r := RoverData{
		Id: r1.Id,
	}
	r3, err := server.Rover(r)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}
