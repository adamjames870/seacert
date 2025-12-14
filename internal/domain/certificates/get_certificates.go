package certificates

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/google/uuid"
)

func GetCertificates(state *internal.ApiState, ctx context.Context) ([]Certificate, error) {

	certTypeMap, errCertTypeMap := getMapOfCertTypes(state, ctx)
	if errCertTypeMap != nil {
		return nil, errCertTypeMap
	}

	issuerMap, errIssuerMap := getMapOfIssuersToCertTypes(state, ctx)
	if errIssuerMap != nil {
		return nil, errIssuerMap
	}

	certs, errCerts := state.Queries.GetCerts(ctx)
	if errCerts != nil {
		return nil, errCerts
	}

	apiCerts := make([]Certificate, 0, len(certs))
	for _, cert := range certs {
		certType, _ := certTypeMap[cert.CertTypeID]
		issuer, _ := issuerMap[cert.IssuerID]
		apiCerts = append(apiCerts, MapCertificateDbToDomain(cert, certType, issuer))
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

	issuer, errIssuer := issuers.GetIssuerFromId(state, ctx, cert.IssuerID.String())
	if errIssuer != nil {
		return Certificate{}, errIssuer
	}

	return MapCertificateDbToDomain(cert, certType, issuer), nil

}

func getMapOfCertTypes(state *internal.ApiState, ctx context.Context) (map[uuid.UUID]cert_types.CertificateType, error) {

	certTypes, errCertTypes := cert_types.GetCertTypes(state, ctx)
	if errCertTypes != nil {
		return nil, errCertTypes
	}

	certTypeMap := make(map[uuid.UUID]cert_types.CertificateType)
	for _, cType := range certTypes {
		certTypeMap[cType.Id] = cType
	}

	return certTypeMap, nil

}

func getMapOfIssuersToCertTypes(state *internal.ApiState, ctx context.Context) (map[uuid.UUID]issuers.Issuer, error) {

	dbIssuers, errIssuers := issuers.GetIssuers(state, ctx)
	if errIssuers != nil {
		return nil, errIssuers
	}

	issuerMap := make(map[uuid.UUID]issuers.Issuer)
	for _, dbIssuer := range dbIssuers {
		issuerMap[dbIssuer.Id] = dbIssuer
	}

	return issuerMap, nil

}
