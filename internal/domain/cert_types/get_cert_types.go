package cert_types

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/google/uuid"
)

func GetCertTypes(state *internal.ApiState, ctx context.Context) ([]CertificateType, error) {

	certTypes, errCertTypes := state.Queries.GetCertTypes(ctx)
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
