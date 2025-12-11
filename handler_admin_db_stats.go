package main

import "net/http"

type responseDbStats struct {
	CountCert     int `json:"countCert"`
	CountCertType int `json:"countCertType"`
}

func (state *apiState) handlerAdminDbStats(w http.ResponseWriter, r *http.Request) {
	if !state.isDev {
		respondWithError(w, 403, "")
	}

	countCert, errCountCert := state.db.CountCertificates(r.Context())
	if errCountCert != nil {
		respondWithError(w, 500, errCountCert.Error())
	}

	countCertType, errCountCertType := state.db.CountCertTypes(r.Context())
	if errCountCertType != nil {
		respondWithError(w, 500, errCountCertType.Error())
	}

	rv := responseDbStats{
		CountCert:     int(countCert),
		CountCertType: int(countCertType),
	}

	respondWithJSON(w, 200, rv)

}
