package notifications

import (
	"net/http"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/notifications"
)

func HandlerNotifyTestGenerate7Day(state *internal.ApiState) http.HandlerFunc {
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

func HandlerNotifyTestGenerate1Month(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(w, r, http.StatusForbidden, "Forbidden", nil)
			return
		}

		generator := notifications.NewGenerator(state.Repo)

		count, err := generator.GenerateNoCertificates1Month(r.Context())
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

func HandlerNotifyTestGenerateExpiring(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(w, r, http.StatusForbidden, "Forbidden", nil)
			return
		}

		generator := notifications.NewGenerator(state.Repo)

		count, err := generator.GenerateCertificateExpirySummary(r.Context(), time.Now())
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

func HandlerNotifyTestSend(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(
				w,
				r,
				http.StatusForbidden,
				"Forbidden",
				nil,
			)
			return
		}

		processor := notifications.NewProcessor(state.Repo)

		result, err := processor.ProcessPendingNotifications(r.Context(), 10)
		if err != nil {
			handlers.RespondWithError(
				w,
				r,
				http.StatusInternalServerError,
				"Error processing pending notifications",
				err,
			)
			return
		}

		handlers.RespondWithJSON(
			w,
			http.StatusOK,
			result,
		)
	}
}
