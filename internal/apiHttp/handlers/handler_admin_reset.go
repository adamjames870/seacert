package handlers

import (
	"fmt"
	"net/http"

	"github.com/adamjames870/seacert/internal"
)

func HandlerApiReset(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !state.IsDev {
			respondWithError(w, 403, "")
		}

		errResetCerts := state.Queries.ResetCerts(r.Context())
		if errResetCerts != nil {
			fmt.Println(errResetCerts)
			respondWithError(w, 500, errResetCerts.Error())
		}

		errResetCertTypes := state.Queries.ResetCertTypes(r.Context())
		if errResetCertTypes != nil {
			fmt.Println(errResetCertTypes)
			respondWithError(w, 500, errResetCertTypes.Error())
		}

		errResetIssuers := state.Queries.ResetIssuers(r.Context())
		if errResetIssuers != nil {
			fmt.Println(errResetIssuers)
			respondWithError(w, 500, errResetIssuers.Error())
		}

		respondWithJSON(w, 200, map[string]string{"message": "db reset"})

	}
}
