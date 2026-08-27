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

		delivery, err := processor.SendDelivery(
			r.Context(),
			notification,
		)
		if err != nil {
			handlers.RespondWithError(
				w,
				r,
				http.StatusInternalServerError,
				"Error sending notification email",
				err,
			)
			return
		}

		handlers.RespondWithJSON(
			w,
			http.StatusOK,
			map[string]any{
				"delivery_id":         delivery.ID,
				"notification_id":     delivery.NotificationID,
				"recipient":           delivery.Recipient,
				"provider":            delivery.Provider,
				"provider_message_id": delivery.ProviderMessageID,
				"status":              delivery.Status,
				"attempt":             delivery.Attempt,
				"sent_at":             delivery.SentAt,
			},
		)
	}
}
