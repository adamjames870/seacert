package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddCertType(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// POST api/cert-types

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsAddCertificateType{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
			return
		}

		certType, certTypeErr := cert_types.WriteNewCertType(state, r.Context(), params)
		if certTypeErr != nil {
			respondWithError(w, 500, certTypeErr.Error())
			return
		}

		rv := cert_types.MapCertificateTypeDomainToDto(certType)

		respondWithJSON(w, 201, rv)

	}
}
