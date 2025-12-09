package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/joho/godotenv"
)

func (state *apiState) LoadState() error {

	errEnv := godotenv.Load()
	if errEnv != nil {
		return fmt.Errorf("error loading .env file: %w", errEnv)
	}

	state.mux = http.NewServeMux()
	err := state.loadDb()
	if err != nil {
		return err
	}
	state.setDevFlag()
	return nil
}

func (state *apiState) setDevFlag() {
	platform := os.Getenv("PLATFORM")
	fmt.Printf("platform = %s\n", platform)
	state.isDev = platform == "dev"
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
