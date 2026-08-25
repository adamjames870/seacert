package notifications

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/notifications"
)

func HandlerNotifyTestEmail(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(w, r, http.StatusForbidden, "Forbidden", nil)
			return
		}

		generator := notifications.NewGenerator(state.Repo)

		count, err := generator.GenerateNoCertificates7Day(r.Context())
		if err != nil {
			handlers.RespondWithError(
				w,
				r,
				http.StatusInternalServerError,
				"Error generating notifications",
				err,
			)
			return
		}

		handlers.RespondWithJSON(
			w,
			http.StatusOK,
			map[string]any{
				"generated": count,
			},
		)
	}
}
