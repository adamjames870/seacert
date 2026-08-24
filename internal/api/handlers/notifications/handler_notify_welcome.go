package notifications

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/email"
)

func HandlerNotifyWelcomeEmail(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authUser, ok := auth.UserFromContext(r.Context())
		if !ok {
			handlers.RespondWithError(w, r, 401, "Unauthorized", nil)
			return
		}

		req := email.Request{
			To:   authUser.Email,
			Name: authUser.Forename,
		}

		sendId, err := email.Welcome(req)
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Error sending email", err)
			return
		}

		handlers.RespondWithJSON(w, http.StatusOK, sendId)
	}
}
