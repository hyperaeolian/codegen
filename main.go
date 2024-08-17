package main

import (
    "fmt"
    "log"
    "net/http"

	"codegen/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/", handlers.Root)
	r.HandleFunc("/{language}/generate_sdk", handlers.CodegenHandler).Methods("POST")
	r.HandleFunc("/list_languages", handlers.Meta)

    port := "8080"
    fmt.Printf("Starting server on port %s...\n", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}
