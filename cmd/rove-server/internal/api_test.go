// +build integration

package internal

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mdiluz/rove/pkg/game"
	"github.com/mdiluz/rove/pkg/rove"
	"github.com/stretchr/testify/assert"
)

const (
	defaultAddress = "localhost:8080"
)

var serv = func() rove.Server {
	var address = os.Getenv("ROVE_SERVER_ADDRESS")
	if len(address) == 0 {
		address = defaultAddress
	}
	return rove.Server(address)
}()

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
	_, err := serv.Register(d1)
	assert.NoError(t, err)

	d2 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	_, err = serv.Register(d2)
	assert.NoError(t, err)

	_, err = serv.Register(d1)
	assert.Error(t, err)
}

func TestServer_Command(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	_, err := serv.Register(d1)
	assert.NoError(t, err)

	c := rove.CommandData{
		Commands: []game.Command{
			{
				Command:  game.CommandMove,
				Bearing:  "N",
				Duration: 1,
			},
		},
	}
	_, err = serv.Command(d1.Name, c)
	assert.NoError(t, err)
}

func TestServer_Radar(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	_, err := serv.Register(d1)
	assert.NoError(t, err)

	_, err = serv.Radar(d1.Name)
	assert.NoError(t, err)
}

func TestServer_Rover(t *testing.T) {
	d1 := rove.RegisterData{
		Name: uuid.New().String(),
	}
	_, err := serv.Register(d1)
	assert.NoError(t, err)

	_, err = serv.Rover(d1.Name)
	assert.NoError(t, err)
}
