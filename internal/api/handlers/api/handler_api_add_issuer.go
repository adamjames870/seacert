package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddIssuer(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// POST api/issuers

		params := dto.ParamsAddIssuer{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		dbIssuer, errIssuer := issuers.CreateIssuer(r.Context(), state.Repo, params)
		if errIssuer != nil {
			code, msg := handlers.MapDomainError(errIssuer)
			handlers.RespondWithError(w, r, code, msg, errIssuer)
			return
		}

		state.Logger.Info("Issuer created", "id", dbIssuer.Id, "name", dbIssuer.Name)

		rv := issuers.MapIssuerDomainToDto(dbIssuer)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
