package issuers

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewIssuer(state *internal.ApiState, ctx context.Context, params dto.ParamsAddIssuer) (Issuer, error) {

	country := domain.ToNullStringFromPointer(params.Country)
	website := domain.ToNullStringFromPointer(params.Website)

	newIssuer := sqlc.CreateIssuerParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Country:   country,
		Website:   website,
	}

	dbIssuer, errWriteNewIssuer := state.Queries.CreateIssuer(ctx, newIssuer)
	if errWriteNewIssuer != nil {
		return Issuer{}, errWriteNewIssuer
	}

	apiIssuer := MapIssuerDbToDomain(dbIssuer)

	return apiIssuer, nil

}
