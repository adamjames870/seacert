package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// POST api/certificates
		params := dto.ParamsAddCertificate{}
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

		cert, certErr := certificates.CreateCertificate(r.Context(), state.Repo, params, userId)
		if certErr != nil {
			code, msg := handlers.MapDomainError(certErr)
			handlers.RespondWithError(w, r, code, msg, certErr)
			return
		}

		state.Logger.Info("Certificate created", "user_id", userId, "certificate_id", cert.Id)
		rv := certificates.MapCertificateDomainToDto(r.Context(), state.Storage, cert)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
