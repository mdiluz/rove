package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/mdiluz/rove/cmd/rove-accountant/internal"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/persistence"
	"google.golang.org/grpc"
)

var data = os.Getenv("DATA_PATH")

// accountantServer is the internal object to manage the requests
type accountantServer struct {
	accountant *internal.Accountant
	sync       sync.RWMutex
}

// Register will register an account
func (a *accountantServer) Register(ctx context.Context, in *accounts.RegisterInfo) (*accounts.RegisterResponse, error) {
	a.sync.Lock()
	defer a.sync.Unlock()

	// Try and register the account itself
	log.Printf("Registering account: %s\n", in.Name)
	if _, err := a.accountant.RegisterAccount(in.Name); err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	// Save out the accounts
	if err := persistence.Save("accounts", a.accountant); err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	return &accounts.RegisterResponse{}, nil
}

// AssignData assigns a key value pair to an account
func (a *accountantServer) AssignValue(_ context.Context, in *accounts.DataKeyValue) (*accounts.DataKeyResponse, error) {
	a.sync.RLock()
	defer a.sync.RUnlock()

	// Try and assign the data
	log.Printf("Assigning value for account %s: %s->%s\n", in.Account, in.Key, in.Value)
	err := a.accountant.AssignData(in.Account, in.Key, in.Value)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	return &accounts.DataKeyResponse{}, nil

}

// GetData gets the value for a key
func (a *accountantServer) GetValue(_ context.Context, in *accounts.DataKey) (*accounts.DataResponse, error) {
	a.sync.RLock()
	defer a.sync.RUnlock()

	// Try and fetch the value
	data, err := a.accountant.GetValue(in.Account, in.Key)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	return &accounts.DataResponse{Value: data}, nil

}

// main
func main() {
	// Get the port
	var iport int
	var port = os.Getenv("PORT")
	if len(port) == 0 {
		iport = 9091
	} else {
		var err error
		iport, err = strconv.Atoi(port)
		if err != nil {
			log.Fatal("$PORT not valid int")
		}
	}

	persistence.SetPath(data)

	// Initialise and load the accountant
	accountant := internal.NewAccountant()
	if err := persistence.Load("accounts", accountant); err != nil {
		log.Fatalf("failed to load account data: %s", err)
	}

	// Set up the RPC server and register
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", iport))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	accounts.RegisterAccountantServer(grpcServer, &accountantServer{
		accountant: accountant,
	})

	// Set up the close handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Quit requested, exiting...")
		grpcServer.Stop()
	}()

	// Serve the RPC server
	log.Printf("Serving accountant on %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %s", err)
	}

	// Save out the accountant data
	if err := persistence.Save("accounts", accountant); err != nil {
		log.Fatalf("failed to save accounts: %s", err)
	}
}
