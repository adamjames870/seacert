package main

import (
	"fmt"
	"net/http"
)

func (state *apiState) handlerApiReset(w http.ResponseWriter, r *http.Request) {
	if !state.isDev {
		respondWithError(w, 403, "")
	}

	errResetCerts := state.db.ResetCerts(r.Context())
	if errResetCerts != nil {
		fmt.Println(errResetCerts)
		respondWithError(w, 500, errResetCerts.Error())
	}

	errResetCertTypes := state.db.ResetCertTypes(r.Context())
	if errResetCertTypes != nil {
		fmt.Println(errResetCertTypes)
		respondWithError(w, 500, errResetCertTypes.Error())
	}

	respondWithJSON(w, 200, map[string]string{"message": "db reset"})

}
