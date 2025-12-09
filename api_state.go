package main

import (
	"net/http"

	"github.com/adamjames870/seacert/internal/database"
)

type apiState struct {
	mux   *http.ServeMux
	db    *database.Queries
	isDev bool
}
