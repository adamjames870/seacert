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
	"github.com/adamjames870/seacert/internal/database/migrations"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/logging"
	"github.com/adamjames870/seacert/internal/repository/postgres"
	"github.com/adamjames870/seacert/internal/storage"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"google.golang.org/api/option"
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
	errGemini := loadGemini(state)
	if errGemini != nil {
		return errGemini
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

func loadGemini(state *internal.ApiState) error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		state.Logger.Warn("GEMINI_API_KEY is not set, certificate extraction will not work")
		return nil
	}

	modelName := os.Getenv("GEMINI_MODEL_NAME")
	if modelName == "" {
		modelName = "gemini-2.0-flash" // Default fallback
		state.Logger.Warn("GEMINI_MODEL_NAME is not set, using default", "default", modelName)
	}
	state.GeminiModelName = modelName

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	state.Gemini = client
	return nil
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

	// Run Migrations
	if err := runMigrations(state, db); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	state.Queries = sqlc.New(db)
	state.Repo = postgres.NewRepository(db)

	return nil
}

func runMigrations(state *internal.ApiState, db *sql.DB) error {
	state.Logger.Info("Running database migrations...")

	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(slog.NewLogLogger(state.Logger.Handler(), slog.LevelInfo))

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Apply migrations
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	state.Logger.Info("Database migrations completed successfully")
	return nil
}
