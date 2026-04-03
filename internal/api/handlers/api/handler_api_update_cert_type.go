package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerUpdateCertType(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// PUT /api/cert-types?id=<uuid>

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			handlers.RespondWithError(w, r, 400, "Missing required parameter 'id'", nil)
			return
		}

		params := dto.ParamsUpdateCertificateType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		params.Id = idParam

		certType, errCertType := cert_types.UpdateCertificateType(r.Context(), state.Repo, params)
		if errCertType != nil {
			code, msg := handlers.MapDomainError(errCertType)
			handlers.RespondWithError(w, r, code, msg, errCertType)
			return
		}

		state.Logger.Info("Certificate type updated", "id", certType.Id, "name", certType.Name)

		certTypeDto := cert_types.MapCertificateTypeDomainToDto(certType)
		handlers.RespondWithJSON(w, 200, certTypeDto)

	}
}
