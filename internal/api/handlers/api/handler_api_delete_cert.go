package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
)

func HandlerApiDeleteCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// DELETE api/certificates?id=<uuid>

		idParam := r.URL.Query().Get("id")

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, 401, "Unauthorized", errId)
			return
		}

		err := certificates.DeleteCertificate(state, r.Context(), idParam, userId)
		if err != nil {
			handlers.RespondWithError(w, 500, "Error deleting certificate", err)
			return
		}

		handlers.RespondWithJSON(w, 200, "certificate deleted")

	}
}
