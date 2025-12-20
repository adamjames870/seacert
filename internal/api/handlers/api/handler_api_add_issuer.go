package api

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
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
			handlers.RespondWithError(w, 400, "Invalid request payload", errDecode)
			return
		}

		dbIssuer, errIssuer := issuers.WriteNewIssuer(state, r.Context(), params)
		if errIssuer != nil {
			handlers.RespondWithError(w, 500, "Error creating issuer", errIssuer)
			return
		}

		state.Logger.Info("Issuer created", "id", dbIssuer.Id, "name", dbIssuer.Name)

		rv := issuers.MapIssuerDomainToDto(dbIssuer)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
