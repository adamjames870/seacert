package certificates

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/google/uuid"
)

func GetCertificates(state *internal.ApiState, ctx context.Context) ([]Certificate, error) {

	certTypes, errCertTypes := cert_types.GetCertTypes(state, ctx)
	if errCertTypes != nil {
		return nil, errCertTypes
	}

	certTypeMap := make(map[uuid.UUID]cert_types.CertificateType)
	for _, cType := range certTypes {
		certTypeMap[cType.Id] = cType
	}

	certs, errCerts := state.Queries.GetCerts(ctx)
	if errCerts != nil {
		return nil, errCerts
	}

	apiCerts := make([]Certificate, 0, len(certs))
	for _, cert := range certs {
		certType, _ := certTypeMap[cert.CertTypeID]
		apiCerts = append(apiCerts, MapCertificateDbToDomain(cert, certType))
	}

	return apiCerts, nil

}

func GetCertificateFromId(state *internal.ApiState, ctx context.Context, id string) (Certificate, error) {

	certUuid, errId := uuid.Parse(id)
	if errId != nil {
		return Certificate{}, errId
	}

	cert, errCert := state.Queries.GetCertFromId(ctx, certUuid)
	if errCert != nil {
		return Certificate{}, errCert
	}

	certType, errCertType := cert_types.GetCertTypeFromId(state, ctx, cert.CertTypeID.String())
	if errCertType != nil {
		return Certificate{}, errCertType
	}

	return MapCertificateDbToDomain(cert, certType), nil

}
