package api

import (
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
			handlers.RespondWithError(w, r, 400, "Missing required parameter 'id'", nil)
			return
		}

		params := dto.ParamsUpdateIssuer{}
		params.Id = idParam

		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		issuer, errIssuer := issuers.UpdateIssuer(r.Context(), state.Repo, params)
		if errIssuer != nil {
			code, msg := handlers.MapDomainError(errIssuer)
			handlers.RespondWithError(w, r, code, msg, errIssuer)
			return
		}

		state.Logger.Info("Issuer updated", "id", issuer.Id, "name", issuer.Name)
		issuerDto := issuers.MapIssuerDomainToDto(issuer)
		handlers.RespondWithJSON(w, 200, issuerDto)

	}
}
