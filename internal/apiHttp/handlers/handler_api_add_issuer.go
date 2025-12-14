package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddIssuer(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// POST api/issuers

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsAddIssuer{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
			return
		}

		dbIssuer, errIssuer := issuers.WriteNewIssuer(state, r.Context(), params)
		if errIssuer != nil {
			respondWithError(w, 500, errIssuer.Error())
			return
		}

		rv := issuers.MapIssuerDomainToDto(dbIssuer)

		respondWithJSON(w, 201, rv)

	}
}
