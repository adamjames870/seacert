package issuers

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/google/uuid"
)

func GetIssuers(state *internal.ApiState, ctx context.Context) ([]Issuer, error) {

	issuers, errIssuers := state.Queries.GetIssuers(ctx)
	if errIssuers != nil {
		return nil, errIssuers
	}

	apiIssuers := make([]Issuer, 0, len(issuers))
	for _, issuer := range issuers {
		apiIssuers = append(apiIssuers, MapIssuerDbToDomain(issuer))
	}

	return apiIssuers, nil

}

func GetIssuerFromId(state *internal.ApiState, ctx context.Context, id string) (Issuer, error) {

	uuidId, errId := uuid.Parse(id)
	if errId != nil {
		return Issuer{}, errId
	}

	issuer, errIssuer := state.Queries.GetIssuerById(ctx, uuidId)
	if errIssuer != nil {
		return Issuer{}, errIssuer
	}

	return MapIssuerDbToDomain(issuer), nil

}

func GetIssuerFromName(state *internal.ApiState, ctx context.Context, name string) (Issuer, error) {

	issuer, errIssuer := state.Queries.GetIssuerByName(ctx, name)
	if errIssuer != nil {
		return Issuer{}, errIssuer
	}

	return MapIssuerDbToDomain(issuer), nil

}
