package api

import (
	"encoding/json"
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
			handlers.RespondWithError(w, 400, "Missing required parameter 'id'", nil)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsUpdateCertificateType{}

		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, 400, "Invalid request payload", errDecode)
			return
		}

		params.Id = idParam

		certType, errCertType := cert_types.UpdateCertificateType(state, r.Context(), params)
		if errCertType != nil {
			handlers.RespondWithError(w, 500, "Error updating certificate type", errCertType)
			return
		}

		state.Logger.Info("Certificate type updated", "id", certType.Id, "name", certType.Name)

		certTypeDto := cert_types.MapCertificateTypeDomainToDto(certType)
		handlers.RespondWithJSON(w, 200, certTypeDto)

	}
}
