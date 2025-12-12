package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/joho/godotenv"
)

func LoadState(state *internal.ApiState) error {

	errEnv := godotenv.Load()
	if errEnv != nil {
		return fmt.Errorf("error loading .env file: %w", errEnv)
	}

	err := loadDb(state)
	if err != nil {
		return err
	}
	setDevFlag(state)
	return nil
}

func setDevFlag(state *internal.ApiState) {
	platform := os.Getenv("PLATFORM")
	fmt.Printf("platform = %s\n", platform)
	state.IsDev = platform == "dev"
}

func loadDb(state *internal.ApiState) error {
	dbUrl := os.Getenv("DB_URL")
	db, errDb := sql.Open("postgres", dbUrl)
	if errDb != nil {
		return fmt.Errorf("unable to load DB: %w", errDb)
	}

	state.Queries = sqlc.New(db)

	return nil
}
