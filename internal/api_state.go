package internal

import (
	"log/slog"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/storage"
)

type ApiState struct {
	Queries *sqlc.Queries
	Storage storage.Storage
	IsDev   bool
	Logger  *slog.Logger
}
