package admin

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerAdminResolveCertType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POST /admin/cert-types/resolve

		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsResolveCertificateType{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, r, 400, "Invalid request payload", errDecode)
			return
		}

		if params.ProvisionalId == "" || params.ReplacementId == "" {
			handlers.RespondWithError(w, r, 400, "Both provisional-id and replacement-id are required", nil)
			return
		}

		errResolve := cert_types.ResolveProvisionalCertType(state, r.Context(), params)
		if errResolve != nil {
			handlers.RespondWithError(w, r, 500, "Error resolving certificate type", errResolve)
			return
		}

		state.Logger.Info("Provisional certificate type resolved", "provisional_id", params.ProvisionalId, "replacement_id", params.ReplacementId)

		handlers.RespondWithJSON(w, 200, "Certificate type resolved and provisional type deleted")
	}
}
