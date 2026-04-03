package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiUpdateCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// PUT api/certificates

		params := dto.ParamsUpdateCertificate{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		params.UserId = userId.String()

		cert, certErr := certificates.UpdateCertificate(r.Context(), state.Repo, state.Storage, state.Logger, params, userId)
		if certErr != nil {
			code, msg := handlers.MapDomainError(certErr)
			handlers.RespondWithError(w, r, code, msg, certErr)
			return
		}

		state.Logger.Info("Certificate updated", "user_id", userId, "certificate_id", cert.Id)
		certDto := certificates.MapCertificateDomainToDto(r.Context(), state.Storage, cert)

		handlers.RespondWithJSON(w, 200, certDto)

	}
}
