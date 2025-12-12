package internal

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
)

type ApiState struct {
	Queries *sqlc.Queries
	IsDev   bool
}
