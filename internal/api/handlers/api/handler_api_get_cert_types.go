package api

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func HandlerApiGetCertTypes(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET /api/cert-types
		// GET /api/cert-types?id=<uuid>
		// GET /api/cert-types?name=<name>

		idParam := r.URL.Query().Get("id")
		nameParam := r.URL.Query().Get("name")

		if idParam == "" && nameParam == "" {
			user, _ := auth.UserFromContext(r.Context())
			isAdmin := user.Role == "admin"
			var userId *uuid.UUID
			if uid, err := uuid.Parse(user.Id); err == nil {
				userId = &uid
			}

			rv, err := getAllCertTypes(r.Context(), state.Repo, userId, isAdmin)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if idParam != "" {
			uid, errParse := uuid.Parse(idParam)
			if errParse != nil {
				handlers.RespondWithError(w, r, 400, "Invalid ID", errParse)
				return
			}
			rv, err := getCertTypeFromId(r.Context(), state.Repo, uid)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if nameParam != "" {
			rv, err := getCertTypeFromName(r.Context(), state.Repo, nameParam)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
		}

	}
}

func getAllCertTypes(ctx context.Context, repo domain.Repository, userId *uuid.UUID, isAdmin bool) ([]dto.CertificateType, error) {
	certTypes, errCertTypes := cert_types.GetCertTypes(ctx, repo, userId, isAdmin)
	if errCertTypes != nil {
		return nil, errCertTypes
	}

	rv := make([]dto.CertificateType, 0, len(certTypes))
	for _, certType := range certTypes {
		rv = append(rv, cert_types.MapCertificateTypeDomainToDto(certType))
	}
	return rv, nil
}

func getCertTypeFromId(ctx context.Context, repo domain.Repository, id uuid.UUID) (dto.CertificateType, error) {
	certType, errCertType := cert_types.GetCertTypeFromId(ctx, repo, id)
	if errCertType != nil {
		return dto.CertificateType{}, errCertType
	}
	return cert_types.MapCertificateTypeDomainToDto(certType), nil
}

func getCertTypeFromName(ctx context.Context, repo domain.Repository, name string) (dto.CertificateType, error) {
	certType, errCertType := cert_types.GetCertTypeFromName(ctx, repo, name)
	if errCertType != nil {
		return dto.CertificateType{}, errCertType
	}
	return cert_types.MapCertificateTypeDomainToDto(certType), nil
}
