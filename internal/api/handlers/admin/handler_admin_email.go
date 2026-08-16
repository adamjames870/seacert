package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/email"
)

func HandlerAdminTestEmail(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !state.IsDev {
			handlers.RespondWithError(w, r, http.StatusForbidden, "Forbidden", nil)
			return
		}

		authUser, ok := auth.UserFromContext(r.Context())
		if !ok {
			handlers.RespondWithError(w, r, 401, "Unauthorized", nil)
			return
		}

		req := email.EmailRequest{
			To:   authUser.Email,
			Name: authUser.Forename,
		}

		email.TestEmail(state, req, w, r)
	}
}
