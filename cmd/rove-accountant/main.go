package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mdiluz/rove/cmd/rove-accountant/internal"
	"github.com/mdiluz/rove/pkg/accounts"
	"github.com/mdiluz/rove/pkg/persistence"
	"google.golang.org/grpc"
)

var address = os.Getenv("HOST_ADDRESS")
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
	fmt.Printf("Registering account: %s\n", in.Name)
	if _, err := a.accountant.RegisterAccount(in.Name); err != nil {
		fmt.Printf("Error: %s\n", err)
		return &accounts.RegisterResponse{Success: false, Error: fmt.Sprintf("error registering account: %s", err)}, nil
	}

	// Save out the accounts
	if err := persistence.Save("accounts", a.accountant); err != nil {
		fmt.Printf("Error: %s\n", err)
		return &accounts.RegisterResponse{Success: false, Error: fmt.Sprintf("failed to save accounts: %s", err)}, nil
	}

	return &accounts.RegisterResponse{Success: true}, nil
}

// AssignData assigns a key value pair to an account
func (a *accountantServer) AssignValue(_ context.Context, in *accounts.DataKeyValue) (*accounts.Response, error) {
	a.sync.RLock()
	defer a.sync.RUnlock()

	// Try and assign the data
	fmt.Printf("Assigning value for account %s: %s->%s\n", in.Account, in.Key, in.Value)
	err := a.accountant.AssignData(in.Account, in.Key, in.Value)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return &accounts.Response{Success: false, Error: err.Error()}, nil
	}

	return &accounts.Response{Success: true}, nil

}

// GetData gets the value for a key
func (a *accountantServer) GetValue(_ context.Context, in *accounts.DataKey) (*accounts.DataResponse, error) {
	a.sync.RLock()
	defer a.sync.RUnlock()

	// Try and fetch the rover
	fmt.Printf("Getting value for account %s: %s\n", in.Account, in.Key)
	data, err := a.accountant.GetValue(in.Account, in.Key)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return &accounts.DataResponse{Success: false, Error: err.Error()}, nil
	}

	return &accounts.DataResponse{Success: true, Value: data}, nil

}

// main
func main() {
	// Verify the input
	if len(address) == 0 {
		log.Fatal("No address set with $HOST_ADDRESS")
	}

	persistence.SetPath(data)

	// Initialise and load the accountant
	accountant := internal.NewAccountant()
	if err := persistence.Load("accounts", accountant); err != nil {
		log.Fatalf("failed to load account data: %s", err)
	}

	// Set up the RPC server and register
	lis, err := net.Listen("tcp", address)
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
		fmt.Println("Quit requested, exiting...")
		grpcServer.Stop()
	}()

	// Serve the RPC server
	fmt.Printf("Serving accountant on %s\n", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server gRPC: %s", err)
	}

	// Save out the accountant data
	if err := persistence.Save("accounts", accountant); err != nil {
		log.Fatalf("failed to save accounts: %s", err)
	}
}
