package notifications

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/notifications"
)

func HandlerNotifyTestGenerate(state *internal.ApiState) http.HandlerFunc {
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

		pending, err := state.Repo.GetPendingNotifications(r.Context())
		if err != nil {
			handlers.RespondWithError(
				w,
				r,
				http.StatusInternalServerError,
				"Error getting pending notifications",
				err,
			)
			return
		}

		if len(pending) == 0 {
			handlers.RespondWithJSON(
				w,
				http.StatusOK,
				map[string]any{
					"message": "No pending notifications found",
				},
			)
			return
		}

		notification := pending[0]

		processor := notifications.NewProcessor(state.Repo)

		err = processor.ProcessNotification(
			r.Context(),
			notification,
		)
		if err != nil {
			handlers.RespondWithError(
				w,
				r,
				http.StatusInternalServerError,
				"Error processing notification",
				err,
			)
			return
		}

		handlers.RespondWithJSON(
			w,
			http.StatusOK,
			map[string]any{
				"message":         "Notification processed successfully",
				"notification_id": notification.ID,
			},
		)
	}
}
