package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	var state apiState

	err := run(&state)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Hello, world!")

}

func run(state *apiState) error {
	errState := state.LoadState()
	if errState != nil {
		return fmt.Errorf("error loading initial api state: %w", errState)
	}

	errEndpoints := state.CreateEndpoints()
	if errEndpoints != nil {
		return fmt.Errorf("error creating endpoints: %w", errEndpoints)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Handler: state.mux,
		Addr:    ":" + port,
	}

	errServe := server.ListenAndServe()
	if errServe != nil {
		return fmt.Errorf("error starting server: %w", errServe)
	}

	return nil
}
