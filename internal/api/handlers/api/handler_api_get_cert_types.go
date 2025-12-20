package api

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiGetCertTypes(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET /api/cert-types
		// GET /api/cert-types?id=<uuid>
		// GET /api/cert-types?name=<name>

		idParam := r.URL.Query().Get("id")
		nameParam := r.URL.Query().Get("name")

		if idParam == "" && nameParam == "" {
			rv, err := getAllCertTypes(state, r.Context())
			if err != nil {
				handlers.RespondWithError(w, 500, "Error fetching certificate types", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if idParam != "" {
			rv, err := getCertTypeFromId(state, r.Context(), idParam)
			if err != nil {
				handlers.RespondWithError(w, 404, "Certificate type not found", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if nameParam != "" {
			rv, err := getCertTypeFromName(state, r.Context(), nameParam)
			if err != nil {
				handlers.RespondWithError(w, 404, "Certificate type not found", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
		}

	}
}

func getAllCertTypes(state *internal.ApiState, ctx context.Context) ([]dto.CertificateType, error) {
	certTypes, errCertTypes := cert_types.GetCertTypes(state, ctx)
	if errCertTypes != nil {
		return nil, errCertTypes
	}

	rv := make([]dto.CertificateType, 0, len(certTypes))
	for _, certType := range certTypes {
		rv = append(rv, cert_types.MapCertificateTypeDomainToDto(certType))
	}
	return rv, nil
}

func getCertTypeFromId(state *internal.ApiState, ctx context.Context, id string) (dto.CertificateType, error) {
	certType, errCertType := cert_types.GetCertTypeFromId(state, ctx, id)
	if errCertType != nil {
		return dto.CertificateType{}, errCertType
	}
	return cert_types.MapCertificateTypeDomainToDto(certType), nil
}

func getCertTypeFromName(state *internal.ApiState, ctx context.Context, name string) (dto.CertificateType, error) {
	certType, errCertType := cert_types.GetCertTypeFromName(state, ctx, name)
	if errCertType != nil {
		return dto.CertificateType{}, errCertType
	}
	return cert_types.MapCertificateTypeDomainToDto(certType), nil
}
