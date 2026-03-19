package api

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func HandlerApiAddCertType(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// POST api/cert-types

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsAddCertificateType{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, 400, "Invalid request payload", errDecode)
			return
		}

		user, _ := auth.UserFromContext(r.Context())
		isAdmin := user.Role == "admin"
		creatorId, _ := uuid.Parse(user.Id)

		certType, errCertType := cert_types.WriteNewCertType(state, r.Context(), params, creatorId, isAdmin)
		if errCertType != nil {
			handlers.RespondWithError(w, 500, "Error creating certificate type", errCertType)
			return
		}

		state.Logger.Info("Certificate type created", "id", certType.Id, "name", certType.Name)

		rv := cert_types.MapCertificateTypeDomainToDto(certType)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
