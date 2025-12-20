package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/logging"
	"github.com/joho/godotenv"
)

func LoadState(state *internal.ApiState) error {

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Printf("Warning: error loading .env file: %v", errEnv)
	}

	state.Logger = logging.NewLogger()
	slog.SetDefault(state.Logger)

	err := loadDb(state)
	if err != nil {
		return err
	}
	setDevFlag(state)
	return nil
}

func setDevFlag(state *internal.ApiState) {
	platform := os.Getenv("PLATFORM")
	log.Printf("platform = %s", platform)
	state.IsDev = platform == "dev"
}

func loadDb(state *internal.ApiState) error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return fmt.Errorf("DB_URL environment variable is not set")
	}

	db, errDb := sql.Open("postgres", dbUrl)
	if errDb != nil {
		return fmt.Errorf("unable to load DB: %w", errDb)
	}

	// Set connection pooling limits
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	state.Queries = sqlc.New(db)

	return nil
}
