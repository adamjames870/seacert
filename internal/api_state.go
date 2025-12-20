package internal

import (
	"log/slog"

	"github.com/adamjames870/seacert/internal/database/sqlc"
)

type ApiState struct {
	Queries *sqlc.Queries
	IsDev   bool
	Logger  *slog.Logger
}
