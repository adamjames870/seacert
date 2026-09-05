package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/repository/postgres"
	"github.com/joho/godotenv"
)

func LoadState(state *notificationState) error {

	_ = godotenv.Load()

	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		log.Fatal("RESEND_API_KEY is not set")
	}

	err := loadDb(state)
	if err != nil {
		return err
	}

	return nil

}

func loadDb(state *notificationState) error {

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
	state.Repo = postgres.NewRepository(db)

	return nil
}
