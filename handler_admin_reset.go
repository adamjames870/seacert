package main

import (
	"fmt"
	"net/http"
)

func (state *apiState) handlerApiReset(w http.ResponseWriter, r *http.Request) {
	if !state.isDev {
		respondWithError(w, 403, "")
	}

	errReset := state.db.ResetDb(r.Context())
	if errReset != nil {
		fmt.Println(errReset)
		respondWithError(w, 500, errReset.Error())
	}

	respondWithJSON(w, 200, map[string]string{"message": "db reset"})

}
