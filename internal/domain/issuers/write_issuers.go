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

func UpdateIssuer(state *internal.ApiState, ctx context.Context, params dto.ParamsUpdateIssuer) (Issuer, error) {

	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return Issuer{}, errParse
	}

	name := domain.ToNullStringFromPointer(params.Name)
	country := domain.ToNullStringFromPointer(params.Country)
	website := domain.ToNullStringFromPointer(params.Website)

	updateIssuer := sqlc.UpdateIssuerParams{
		ID:      uuidId,
		Name:    name,
		Country: country,
		Website: website,
	}

	dbIssuer, errUpdateIssuer := state.Queries.UpdateIssuer(ctx, updateIssuer)
	if errUpdateIssuer != nil {
		return Issuer{}, errUpdateIssuer
	}

	apiIssuer := MapIssuerDbToDomain(dbIssuer)
	return apiIssuer, nil

}
