package api

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddCert(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// POST api/certificates
		decoder := json.NewDecoder(r.Body)
		params := dto.ParamsAddCertificate{}
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

		params.UserId = userId.String()

		cert, certErr := certificates.WriteNewCert(state, r.Context(), params)
		if certErr != nil {
			handlers.RespondWithError(w, 500, "error writing cert: "+certErr.Error())
			return
		}

		rv := certificates.MapCertificateDomainToDto(cert)

		handlers.RespondWithJSON(w, 201, rv)

	}
}
