package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
)

func HandlerAdminGetUser(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET admin/users

		authUser, ok := auth.UserFromContext(r.Context())
		if !ok {
			handlers.RespondWithError(w, r, 401, "Unauthorized", nil)
			return
		}

		uuidId, errParse := uuid.Parse(authUser.Id)
		if errParse != nil {
			handlers.RespondWithError(w, r, 500, "Invalid user ID format", errParse)
			return
		}

		apiUser, errUser := users.GetUser(r.Context(), state.Repo, uuidId)
		if errUser != nil {
			code, msg := handlers.MapDomainError(errUser)
			handlers.RespondWithError(w, r, code, msg, errUser)
			return
		}

		apiUser.Role = authUser.Role
		userDto := users.MapUserDomainToDto(apiUser)

		handlers.RespondWithJSON(w, 200, userDto)

	}

}
