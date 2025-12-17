package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/apiHttp/auth"
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
			respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			respondWithError(w, 401, "user not found in context")
			return
		}

		params.Id = userId.String()

		user, userErr := users.UpdateUser(state, r.Context(), params)
		if userErr != nil {
			respondWithError(w, 500, userErr.Error())
		}

		userDto := users.MapUserDomainToDto(user)

		respondWithJSON(w, 200, userDto)

	}

}
