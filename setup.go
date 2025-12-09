package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal/database"
)

func (state *apiState) LoadState() error {
	state.mux = http.NewServeMux()
	err := state.loadDb()
	if err != nil {
		return err
	}
	return nil
}

func (state *apiState) CreateEndpoints() error {

	// ----------- API Handlers ----------------
	state.mux.HandleFunc("GET /api/healthz", healthzHandler)

	return nil
}

func (state *apiState) loadDb() error {
	dbUrl := os.Getenv("DB_URL")
	db, errDb := sql.Open("postgres", dbUrl)
	if errDb != nil {
		return fmt.Errorf("unable to load DB: %w", errDb)
	}

	state.db = database.New(db)

	return nil
}
