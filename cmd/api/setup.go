package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/logging"
	"github.com/adamjames870/seacert/internal/repository/postgres"
	"github.com/adamjames870/seacert/internal/storage"
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
	errStorage := loadStorage(state)
	if errStorage != nil {
		return errStorage
	}
	setDevFlag(state)
	return nil
}

func loadStorage(state *internal.ApiState) error {
	cfg := storage.Config{
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		Endpoint:        os.Getenv("R2_ENDPOINT"),
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
	}

	if cfg.BucketName == "" || cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		state.Logger.Warn("R2 storage configuration is incomplete, uploads will not work")
		return nil
	}

	r2, err := storage.NewR2Storage(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize R2 storage: %w", err)
	}

	state.Storage = r2
	return nil
}

func setDevFlag(state *internal.ApiState) {
	platform := os.Getenv("PLATFORM")
	log.Printf("platform = %s", platform)
	state.IsDev = platform == "dev" || platform == "test"
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
	state.Repo = postgres.NewRepository(db)

	return nil
}
