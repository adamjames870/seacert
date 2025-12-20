package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api"
	_ "github.com/lib/pq"
)

// go run ./cmd/api

func main() {

	var state internal.ApiState

	err := run(&state)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
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
	server := &http.Server{
		Handler:      mux,
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}
