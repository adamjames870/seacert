package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
)

func HandlerAdminReset(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !state.IsDev {
			handlers.RespondWithError(w, 403, "Forbidden", nil)
			return
		}

		errResetCerts := state.Queries.ResetCerts(r.Context())
		if errResetCerts != nil {
			handlers.RespondWithError(w, 500, "Error resetting certificates", errResetCerts)
			return
		}

		errResetCertTypes := state.Queries.ResetCertTypes(r.Context())
		if errResetCertTypes != nil {
			handlers.RespondWithError(w, 500, "Error resetting certificate types", errResetCertTypes)
			return
		}

		errResetIssuers := state.Queries.ResetIssuers(r.Context())
		if errResetIssuers != nil {
			handlers.RespondWithError(w, 500, "Error resetting issuers", errResetIssuers)
			return
		}

		errResetUsers := state.Queries.ResetUsers(r.Context())
		if errResetUsers != nil {
			handlers.RespondWithError(w, 500, "Error resetting users", errResetUsers)
			return
		}

		handlers.RespondWithJSON(w, 200, map[string]string{"message": "db reset"})

	}
}
