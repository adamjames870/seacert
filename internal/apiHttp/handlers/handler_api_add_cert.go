package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// POST api/certificates
		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsAddCertificate{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
			return
		}

		cert, certErr := certificates.WriteNewCert(r.Context(), *state.Queries, params)
		if certErr != nil {
			respondWithError(w, 500, "error writing cert: "+certErr.Error())
			return
		}

		certType, certTypeError := state.Queries.GetCertTypeFromId(r.Context(), cert.CertTypeID)
		if certTypeError != nil {
			respondWithError(w, 500, "error retrieving cert type: "+certTypeError.Error())
		}

		rv := certificates.MapCertificate(cert, certType)

		respondWithJSON(w, 201, rv)

	}
}
