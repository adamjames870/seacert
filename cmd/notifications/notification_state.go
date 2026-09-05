package main

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
)

type notificationState struct {
	Queries *sqlc.Queries
	Repo    domain.Repository
}
