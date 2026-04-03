package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/google/uuid"
)

func HandlerApiDeleteCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// DELETE api/certificates?id=<uuid>

		idParam := r.URL.Query().Get("id")

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		certUuid, errParse := uuid.Parse(idParam)
		if errParse != nil {
			handlers.RespondWithError(w, r, 400, "Invalid certificate ID", errParse)
			return
		}

		err := certificates.DeleteCertificate(r.Context(), state.Repo, state.Storage, state.Logger, certUuid, userId)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		handlers.RespondWithJSON(w, 200, "certificate deleted")

	}
}
