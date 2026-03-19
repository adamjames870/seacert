package cert_types

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

func GetCertTypes(state *internal.ApiState, ctx context.Context, userId *uuid.UUID, isAdmin bool) ([]CertificateType, error) {

	var certTypes []sqlc.CertificateType
	var errCertTypes error

	if isAdmin {
		certTypes, errCertTypes = state.Queries.GetCertTypes(ctx)
	} else if userId != nil {
		certTypes, errCertTypes = state.Queries.GetCertTypesForUser(ctx, uuid.NullUUID{UUID: *userId, Valid: true})
	} else {
		// Should not happen if authenticated, but fallback to approved only
		certTypes, errCertTypes = state.Queries.GetCertTypesForUser(ctx, uuid.NullUUID{Valid: false})
	}

	if errCertTypes != nil {
		return nil, errCertTypes
	}

	apiCertTypes := make([]CertificateType, 0, len(certTypes))
	for _, cType := range certTypes {
		apiCertTypes = append(apiCertTypes, MapCertificateTypeDbToDomain(cType))
	}

	return apiCertTypes, nil

}

func GetCertTypeFromId(state *internal.ApiState, ctx context.Context, id string) (CertificateType, error) {

	uuidId, errId := uuid.Parse(id)
	if errId != nil {
		return CertificateType{}, errId
	}

	certType, errCertType := state.Queries.GetCertTypeFromId(ctx, uuidId)
	if errCertType != nil {
		return CertificateType{}, errCertType
	}

	return MapCertificateTypeDbToDomain(certType), nil

}

func GetCertTypeFromName(state *internal.ApiState, ctx context.Context, name string) (CertificateType, error) {

	certType, errCertType := state.Queries.GetCertTypeFromName(ctx, name)
	if errCertType != nil {
		return CertificateType{}, errCertType
	}

	return MapCertificateTypeDbToDomain(certType), nil

}
