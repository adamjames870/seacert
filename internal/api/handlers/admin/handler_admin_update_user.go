package admin

import (
	"encoding/json"
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

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsUpdateUser{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, 400, "unable to decode json: "+errDecode.Error())
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, 401, "user not found in context")
			return
		}

		params.Id = userId.String()

		user, userErr := users.UpdateUser(state, r.Context(), params)
		if userErr != nil {
			handlers.RespondWithError(w, 500, userErr.Error())
		}

		userDto := users.MapUserDomainToDto(user)

		handlers.RespondWithJSON(w, 200, userDto)

	}

}
