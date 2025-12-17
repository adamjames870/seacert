package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/apiHttp/auth"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
)

func HandlerAdminGetUser(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET admin/users

		authUser, ok := auth.UserFromContext(r.Context())
		if !ok {
			respondWithError(w, 401, "user not found in context")
			return
		}

		uuidId, errParse := uuid.Parse(authUser.Id)
		if errParse != nil {
			respondWithError(w, 500, errParse.Error())
			return
		}

		apiUser, errUser := users.GetUser(state, r.Context(), uuidId)
		if errUser != nil {
			respondWithError(w, 500, errUser.Error())
			return
		}

		userDto := users.MapUserDomainToDto(apiUser)

		respondWithJSON(w, 200, userDto)

	}

}
