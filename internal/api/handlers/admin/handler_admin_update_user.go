package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerAdminUpdateUser(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// PUT admin/users

		params := dto.ParamsUpdateUser{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		params.Id = userId.String()

		user, userErr := users.UpdateUser(r.Context(), state.Repo, params)
		if userErr != nil {
			code, msg := handlers.MapDomainError(userErr)
			handlers.RespondWithError(w, r, code, msg, userErr)
			return
		}

		state.Logger.Info("User updated", "user_id", user.Id)
		userDto := users.MapUserDomainToDto(user)

		handlers.RespondWithJSON(w, 200, userDto)

	}

}
