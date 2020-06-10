package internal

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/stretchr/testify/assert"
)

// To be set by the main function
var serv rove.Server

func TestMain(m *testing.M) {
	s := NewServer()
	if err := s.Initialise(true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serv = rove.Server(s.Addr())

	go s.Run()

	fmt.Printf("Test server hosted on %s", serv)
	code := m.Run()

	if err := s.StopAndClose(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	d2 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r2, err := serv.Register(d2)
	assert.NoError(t, err)
	assert.True(t, r2.Success)

	r3, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.False(t, r3.Success)
}

func TestServer_Command(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	r1, err := serv.Register(d1)
	assert.NoError(t, err)
	assert.True(t, r1.Success)

	c := rove.CommandData{
		Commands: []game.Command{
			{
				Command:  game.CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}
	r3, err := serv.Command(d1.Name, c)
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

	r3, err := serv.Radar(d1.Name)
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

	r3, err := serv.Rover(d1.Name)
	assert.NoError(t, err)
	assert.True(t, r3.Success)
}
