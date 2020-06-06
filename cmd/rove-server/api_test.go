package main

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/mdiluz/rove/pkg/server"
	"github.com/stretchr/testify/assert"
)

var serv rove.Server = "localhost:8080"

func TestMain(m *testing.M) {
	s := server.NewServer(server.OptionPort(8080))
	s.Initialise()
	go s.Run()

	code := m.Run()

	s.Close()

	os.Exit(code)
}

func TestServer_Status(t *testing.T) {
	status, err := serv.Status()
	assert.NoError(t, err)
	assert.True(t, status.Ready)
	assert.NotZero(t, len(status.Version))
}

func TestServer_Register(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	d2 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r2, err := serv.Register(d2)
	assert.NoError(t, err)
	assert.True(t, r2.Success)
	assert.NotZero(t, len(r2.Id))

	r3, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.False(t, r3.Success)
}

func TestServer_Spawn(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := rove.SpawnData{}
	r2, err := serv.Spawn(r1.Id, s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)
}

func TestServer_Command(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := rove.SpawnData{}
	r2, err := serv.Spawn(r1.Id, s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	c := rove.CommandData{
		Commands: []rove.Command{
			{
				Command:  rove.CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}
	r3, err := serv.Command(r1.Id, c)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}

func TestServer_Radar(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := rove.SpawnData{}
	r2, err := serv.Spawn(r1.Id, s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	r3, err := serv.Radar(r1.Id)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}

func TestServer_Rover(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)
	assert.NotZero(t, len(r1.Id))

	s := rove.SpawnData{}
	r2, err := serv.Spawn(r1.Id, s)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	r3, err := serv.Rover(r1.Id)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}
