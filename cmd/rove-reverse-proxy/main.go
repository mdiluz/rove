package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
		endpoint = "localhost:9090"
	}

	var iport int
	var port = os.Getenv("PORT")
	if len(port) == 0 {
		iport = 8080
	} else {
		var err error
		iport, err = strconv.Atoi(port)
		if err != nil {
			log.Fatal("$PORT not valid int")
		}
	}

	// Create a new mux and register it with the gRPC endpoint
	fmt.Printf("Hosting reverse-proxy on %d for %s\n", iport, endpoint)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := rove.RegisterRoveHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server and proxy calls to gRPC endpoint when needed
	if err := http.ListenAndServe(fmt.Sprintf(":%d", iport), mux); err != nil {
		log.Fatal(err)
	}
}
