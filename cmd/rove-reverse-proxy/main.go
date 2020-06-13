package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/mdiluz/rove/pkg/rove"
)

var endpoint = os.Getenv("GRPC_ENDPOINT")
var address = os.Getenv("HOST_ADDRESS")

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a new mux and register it with the gRPC engpoint
	fmt.Printf("Hosting reverse-proxy on %s for %s\n", address, endpoint)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := rove.RegisterRoveHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}

	// Start the HTTP server and proxy calls to gRPC endpoint when needed
	return http.ListenAndServe(address, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
