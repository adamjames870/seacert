package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
)

func HandlerAdminReset(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !state.IsDev {
			handlers.RespondWithError(w, r, 403, "Forbidden", nil)
			return
		}

		err := state.Repo.ResetAll(r.Context())
		if err != nil {
			handlers.RespondWithError(w, r, 500, "Error resetting database", err)
			return
		}

		handlers.RespondWithJSON(w, 200, map[string]string{"message": "db reset"})

	}
}
