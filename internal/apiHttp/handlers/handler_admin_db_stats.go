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

		rv := dto.DbStats{
			CountCert:     int(countCert),
			CountCertType: int(countCertType),
		}

		respondWithJSON(w, 200, rv)
	}

}
