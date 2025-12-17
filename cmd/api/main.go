package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api"
	_ "github.com/lib/pq"
)

// go run ./cmd/api

func main() {

	var state internal.ApiState

	err := run(&state)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Hello, world!")

}

func run(state *internal.ApiState) error {

	errState := LoadState(state)
	if errState != nil {
		return fmt.Errorf("error loading initial api state: %w", errState)
	}

	mux, errEndpoints := api.BuildRouter(state)
	if errEndpoints != nil {
		return fmt.Errorf("error creating endpoints: %w", errEndpoints)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	errServe := server.ListenAndServe()
	if errServe != nil {
		return fmt.Errorf("error starting server: %w", errServe)
	}

	return nil
}
