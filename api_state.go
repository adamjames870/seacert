package main

import "net/http"

type apiState struct {
	mux *http.ServeMux
}
