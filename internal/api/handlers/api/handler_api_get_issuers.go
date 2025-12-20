package api

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiGetIssuers(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET /api/issuers
		// GET /api/issuers?id=<uuid>
		// GET /api/issuers?name=<name>

		idParam := r.URL.Query().Get("id")
		nameParam := r.URL.Query().Get("name")

		if idParam == "" && nameParam == "" {
			rv, err := getAllIssuers(state, r.Context())
			if err != nil {
				handlers.RespondWithError(w, 500, "Error fetching issuers", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if idParam != "" {
			rv, err := getIssuerFromId(state, r.Context(), idParam)
			if err != nil {
				handlers.RespondWithError(w, 404, "Issuer not found", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
			return
		}

		if nameParam != "" {
			rv, err := getIssuerFromName(state, r.Context(), nameParam)
			if err != nil {
				handlers.RespondWithError(w, 404, "Issuer not found", err)
				return
			}
			handlers.RespondWithJSON(w, 200, rv)
		}

	}
}

func getAllIssuers(state *internal.ApiState, ctx context.Context) ([]dto.Issuer, error) {
	dbIssuers, errIssuers := issuers.GetIssuers(state, ctx)
	if errIssuers != nil {
		return nil, errIssuers
	}

	rv := make([]dto.Issuer, 0, len(dbIssuers))
	for _, dbIssuer := range dbIssuers {
		rv = append(rv, issuers.MapIssuerDomainToDto(dbIssuer))
	}
	return rv, nil
}

func getIssuerFromId(state *internal.ApiState, ctx context.Context, id string) (dto.Issuer, error) {
	dbIssuer, erIssuer := issuers.GetIssuerFromId(state, ctx, id)
	if erIssuer != nil {
		return dto.Issuer{}, erIssuer
	}
	return issuers.MapIssuerDomainToDto(dbIssuer), nil
}

func getIssuerFromName(state *internal.ApiState, ctx context.Context, name string) (dto.Issuer, error) {
	dbIssuer, errIssuer := issuers.GetIssuerFromName(state, ctx, name)
	if errIssuer != nil {
		return dto.Issuer{}, errIssuer
	}
	return issuers.MapIssuerDomainToDto(dbIssuer), nil
}
