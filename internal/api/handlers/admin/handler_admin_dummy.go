package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
)

func HandlerAdminDummies(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(w, r, 403, "Forbidden", nil)
			return
		}

		errIssuers := state.Queries.CreateDummyIssuers(r.Context())
		if errIssuers != nil {
			http.Error(w, "Failed to create dummy issuers", http.StatusInternalServerError)
			state.Logger.Info("Failed to create dummy issuers", "error", errIssuers)
			return
		}

		state.Logger.Info("Created dummy issuers")

		errCertTypes := state.Queries.CreateDummyCertTypes(r.Context())
		if errCertTypes != nil {
			http.Error(w, "Failed to create dummy cert types", http.StatusInternalServerError)
			state.Logger.Info("Failed to create dummy cert types", "error", errCertTypes)
			return
		}

		state.Logger.Info("Created dummy cert types")

	}

}
