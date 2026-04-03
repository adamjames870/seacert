package api

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func HandlerApiGetIssuers(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET /api/issuers
		// GET /api/issuers?id=<uuid>
		// GET /api/issuers?name=<name>

		idParam := r.URL.Query().Get("id")
		nameParam := r.URL.Query().Get("name")

		if idParam == "" && nameParam == "" {
			rv, err := getAllIssuers(r.Context(), state.Repo)
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
			rv, err := getIssuerFromId(r.Context(), state.Repo, uid)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if nameParam != "" {
			rv, err := getIssuerFromName(r.Context(), state.Repo, nameParam)
			if err != nil {
				code, msg := handlers.MapDomainError(err)
				handlers.RespondWithError(w, r, code, msg, err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
		}

	}
}

func getAllIssuers(ctx context.Context, repo domain.Repository) ([]dto.Issuer, error) {
	dbIssuers, errIssuers := issuers.GetIssuers(ctx, repo)
	if errIssuers != nil {
		return nil, errIssuers
	}

	rv := make([]dto.Issuer, 0, len(dbIssuers))
	for _, dbIssuer := range dbIssuers {
		rv = append(rv, issuers.MapIssuerDomainToDto(dbIssuer))
	}
	return rv, nil
}

func getIssuerFromId(ctx context.Context, repo domain.Repository, id uuid.UUID) (dto.Issuer, error) {
	dbIssuer, erIssuer := issuers.GetIssuerById(ctx, repo, id)
	if erIssuer != nil {
		return dto.Issuer{}, erIssuer
	}
	return issuers.MapIssuerDomainToDto(dbIssuer), nil
}

func getIssuerFromName(ctx context.Context, repo domain.Repository, name string) (dto.Issuer, error) {
	dbIssuer, errIssuer := issuers.GetIssuerByName(ctx, repo, name)
	if errIssuer != nil {
		return dto.Issuer{}, errIssuer
	}
	return issuers.MapIssuerDomainToDto(dbIssuer), nil
}
