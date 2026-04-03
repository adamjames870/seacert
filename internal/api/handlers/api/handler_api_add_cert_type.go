package api

import (
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

		params := dto.ParamsAddCertificateType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		user, _ := auth.UserFromContext(r.Context())
		isAdmin := user.Role == "admin"
		creatorId, _ := uuid.Parse(user.Id)

		certType, errCertType := cert_types.CreateCertType(r.Context(), state.Repo, params, creatorId, isAdmin)
		if errCertType != nil {
			code, msg := handlers.MapDomainError(errCertType)
			handlers.RespondWithError(w, r, code, msg, errCertType)
			return
		}

		state.Logger.Info("Certificate type created", "id", certType.Id, "name", certType.Name)

		rv := cert_types.MapCertificateTypeDomainToDto(certType)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
