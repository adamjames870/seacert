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

		cert, certErr := certificates.WriteNewCert(state, r.Context(), params)
		if certErr != nil {
			respondWithError(w, 500, "error writing cert: "+certErr.Error())
			return
		}

		rv := certificates.MapCertificateDomainToDto(cert)

		respondWithJSON(w, 201, rv)

	}
}
