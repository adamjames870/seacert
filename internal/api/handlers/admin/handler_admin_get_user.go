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
			handlers.RespondWithError(w, 401, "Unauthorized", nil)
			return
		}

		uuidId, errParse := uuid.Parse(authUser.Id)
		if errParse != nil {
			handlers.RespondWithError(w, 500, "Invalid user ID format", errParse)
			return
		}

		apiUser, errUser := users.GetUser(state, r.Context(), uuidId)
		if errUser != nil {
			handlers.RespondWithError(w, 500, "Error fetching user", errUser)
			return
		}

		userDto := users.MapUserDomainToDto(apiUser)

		handlers.RespondWithJSON(w, 200, userDto)

	}

}
