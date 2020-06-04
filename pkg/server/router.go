package server

import (
	"fmt"
	"io"
	"net/http"
)

// RequestHandler describes a function that handles any incoming request and can respond
type RequestHandler func(io.ReadCloser, io.Writer) error

// Route defines the information for a single path->function route
type Route struct {
	path    string
	method  string
	handler RequestHandler
}

// requestHandlerHTTP wraps a request handler in http checks
func requestHandlerHTTP(method string, handler RequestHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		fmt.Printf("%s\t%s\n", r.Method, r.RequestURI)

		// Verify we're hit with the right method
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)

		} else if err := handler(r.Body, w); err != nil {
			// Log the error
			fmt.Printf("Failed to handle http request: %s", err)

			// Respond that we've had an error
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			// Be a good citizen and set the header for the return
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

		}
	}
}

// NewRouter sets up the server mux
func (s *Server) SetUpRouter() {

	// Array of all our routes
	var routes = []Route{
		{
			path:    "/status",
			method:  http.MethodGet,
			handler: s.HandleStatus,
		},
		{
			path:    "/register",
			method:  http.MethodPost,
			handler: s.HandleRegister,
		},
		{
			path:    "/spawn",
			method:  http.MethodPost,
			handler: s.HandleSpawn,
		},
		{
			path:    "/commands",
			method:  http.MethodPost,
			handler: s.HandleCommands,
		},
		{
			path:    "/view",
			method:  http.MethodPost,
			handler: s.HandleView,
		},
	}

	// Set up the handlers
	for _, route := range routes {
		s.router.HandleFunc(route.path, requestHandlerHTTP(route.method, route.handler))
	}
}
