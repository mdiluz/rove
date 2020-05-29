package rove

type Connection struct {
	address string
}

func NewConnection(address string) *Connection {
	return &Connection{
		address: address,
	}
}

// ServerStatus is a struct that contains information on the status of the server
type ServerStatus struct {
	Ready bool `json:"ready"`
}

func (c *Connection) Status() ServerStatus {
	return ServerStatus{}
}
