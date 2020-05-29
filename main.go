package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", HandleRoot)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
