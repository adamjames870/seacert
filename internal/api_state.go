package internal

import (
	"log/slog"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/storage"
	"github.com/google/generative-ai-go/genai"
)

type ApiState struct {
	Queries         *sqlc.Queries
	Repo            domain.Repository
	Storage         storage.Storage
	IsDev           bool
	Logger          *slog.Logger
	Gemini          *genai.Client
	GeminiModelName string
}
