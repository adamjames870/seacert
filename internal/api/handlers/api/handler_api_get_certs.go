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

		idsParam := r.URL.Query()["id"]

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		switch len(idsParam) {

		case 0:
			// No IDs -> Fetch all certificates
			rv, err := getAllCertificates(state, r.Context(), userId)
			if err != nil {
				handlers.RespondWithError(w, r, 500, "Error fetching certificates", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return

		case 1:
			// 1 ID -> Fetch certificate by ID
			certUuid, errParse := uuid.Parse(idsParam[0])
			if errParse != nil {
				handlers.RespondWithError(w, r, 400, "Invalid certificate ID", errParse)
				return
			}
			rv, err := certificates.GetCertificateById(r.Context(), state.Repo, certUuid, userId)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, certificates.MapCertificateDomainToDto(r.Context(), state.Storage, rv))
			return

		default:
			handlers.RespondWithError(w, r, 400, "Too many ids", nil)
			return

		}

	}
}

func getAllCertificates(state *internal.ApiState, ctx context.Context, userId uuid.UUID) ([]dto.Certificate, error) {

	certs, errCerts := certificates.GetCertificates(ctx, state.Repo, userId)
	if errCerts != nil {
		return nil, errCerts
	}

	rv := make([]dto.Certificate, 0, len(certs))
	for _, cert := range certs {
		rv = append(rv, certificates.MapCertificateDomainToDto(ctx, state.Storage, cert))
	}

	return rv, nil

}
