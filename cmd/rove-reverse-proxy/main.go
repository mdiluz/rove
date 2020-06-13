package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/mdiluz/rove/pkg/rove"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var endpoint = os.Getenv("ROVE_GRPC")
	if len(endpoint) == 0 {
		log.Fatal("Must set $ROVE_GRPC")
	}

	var address = os.Getenv("ROVE_HTTP")
	if len(address) == 0 {
		log.Fatal("Must set $ROVE_HTTP")
	}

	// Create a new mux and register it with the gRPC endpoint
	fmt.Printf("Hosting reverse-proxy on %s for %s\n", address, endpoint)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := rove.RegisterRoveHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server and proxy calls to gRPC endpoint when needed
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
