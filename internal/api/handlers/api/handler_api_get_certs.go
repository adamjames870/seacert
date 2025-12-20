package api

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func HandlerApiGetCerts(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates
		// GET api/certificates?id=<uuid>

		idParam := r.URL.Query().Get("id")

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, 401, "Unauthorized", errId)
			return
		}

		if idParam == "" {
			rv, err := getAllCertificates(state, r.Context(), userId)
			if err != nil {
				handlers.RespondWithError(w, 500, "Error fetching certificates", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if idParam != "" {
			rv, err := certificates.GetCertificateFromId(state, r.Context(), idParam, userId)
			if err != nil {
				handlers.RespondWithError(w, 404, "Certificate not found", err)
				return
			}
			handlers.RespondWithJSON(w, 200, certificates.MapCertificateDomainToDto(rv))
			return
		}

	}
}

func getAllCertificates(state *internal.ApiState, ctx context.Context, userId uuid.UUID) ([]dto.Certificate, error) {

	certs, errCerts := certificates.GetCertificates(state, ctx, userId)
	if errCerts != nil {
		return nil, errCerts
	}

	rv := make([]dto.Certificate, 0, len(certs))
	for _, cert := range certs {
		rv = append(rv, certificates.MapCertificateDomainToDto(cert))
	}

	return rv, nil

}
