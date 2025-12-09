package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	var state apiState
	if state.LoadState() != nil {
		fmt.Println("Error loading initial api state")
		os.Exit(1)
	}

	if godotenv.Load() != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	state.mux = http.NewServeMux()
	if state.CreateEndpoints() != nil {
		fmt.Println("Error creating endpoints")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Handler: state.mux,
		Addr:    ":" + port,
	}
	if server.ListenAndServe() != nil {
		fmt.Println("Error starting server")
		os.Exit(1)
	}

	fmt.Println("Hello, world!")

}
