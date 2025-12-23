package api

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerUpdateIssuer(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// PUT api/issuers?id=<uuid>

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			handlers.RespondWithError(w, 400, "Missing required parameter 'id'", nil)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsUpdateIssuer{}

		params.Id = idParam

		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, 400, "Invalid request payload", errDecode)
			return
		}

		issuer, errIssuer := issuers.UpdateIssuer(state, r.Context(), params)
		if errIssuer != nil {
			handlers.RespondWithError(w, 500, "Error updating issuer", errIssuer)
			return
		}

		state.Logger.Info("Issuer updated", "id", issuer.Id, "name", issuer.Name)
		issuerDto := issuers.MapIssuerDomainToDto(issuer)
		handlers.RespondWithJSON(w, 200, issuerDto)

	}
}
