package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerAdminDbStats(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !state.IsDev {
			respondWithError(w, 403, "")
		}

		countCert, errCountCert := state.Queries.CountCertificates(r.Context())
		if errCountCert != nil {
			respondWithError(w, 500, errCountCert.Error())
		}

		countCertType, errCountCertType := state.Queries.CountCertTypes(r.Context())
		if errCountCertType != nil {
			respondWithError(w, 500, errCountCertType.Error())
		}

		countIssuers, errCountIssuers := state.Queries.CountIssuers(r.Context())
		if errCountIssuers != nil {
			respondWithError(w, 500, errCountIssuers.Error())
		}

		rv := dto.DbStats{
			CountCert:     int(countCert),
			CountCertType: int(countCertType),
			CountIssuer:   int(countIssuers),
		}

		respondWithJSON(w, 200, rv)
	}

}
