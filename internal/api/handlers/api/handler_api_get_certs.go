package api

import (
	"context"
	"database/sql"
	"errors"
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

		idsParam := r.URL.Query()["id"]

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, 401, "Unauthorized", errId)
			return
		}

		switch len(idsParam) {

		case 0:
			// No IDs -> Fetch all certificates
			rv, err := getAllCertificates(state, r.Context(), userId)
			if err != nil {
				handlers.RespondWithError(w, 500, "Error fetching certificates", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return

		case 1:
			// 1 ID -> Fetch certificate by ID
			rv, err := certificates.GetCertificateFromId(state, r.Context(), idsParam[0], userId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					handlers.RespondWithError(w, 404, "Certificate not found", err)
				} else {
					handlers.RespondWithError(w, 500, "Error fetching certificate", err)
				}
				return
			}
			handlers.RespondWithJSON(w, 200, certificates.MapCertificateDomainToDto(rv))
			return

		default:
			handlers.RespondWithError(w, 400, "Too many ids", nil)
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
