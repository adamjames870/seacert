package certificates

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/google/uuid"
)

func GetCertificates(state *internal.ApiState, ctx context.Context, userId uuid.UUID) ([]Certificate, error) {

	uuidId, errParse := uuid.Parse(userId.String())
	if errParse != nil {
		return nil, errParse
	}

	certs, errCerts := state.Queries.GetCerts(ctx, uuidId)
	if errCerts != nil {
		return nil, errCerts
	}

	apiCerts := make([]Certificate, 0, len(certs))
	for _, cert := range certs {
		thisCert := MapCertificateViewDbToDomain(cert.ToCertView())
		thisCert.calculateExpiryDate()
		if cert.HasPredecessors {
			predecessorIds, errPredecessorIds := state.Queries.GetPredecessors(ctx, thisCert.Id)
			if errPredecessorIds != nil {
				return nil, errPredecessorIds
			}
			predecessors, errPredecessors := GetCertificateFromListOfIds(state, ctx, predecessorIds, userId)
			if errPredecessors != nil {
				return nil, errPredecessors
			}
			thisCert.Predecessors = predecessors
		}
		if cert.HasSuccessor {
			successorIds, errSuccessorIds := state.Queries.GetSuccessors(ctx, thisCert.Id)
			if errSuccessorIds != nil {
				return nil, errSuccessorIds
			}
			successors, errSuccessors := GetCertificateFromListOfIds(state, ctx, successorIds, userId)
			if errSuccessors != nil {
				return nil, errSuccessors
			}
			thisCert.Successors = successors
		}
		apiCerts = append(apiCerts, thisCert)
	}

	return apiCerts, nil

}

func GetCertificateFromId(state *internal.ApiState, ctx context.Context, certId string, userId uuid.UUID) (Certificate, error) {

	certUuid, errId := uuid.Parse(certId)
	if errId != nil {
		return Certificate{}, errId
	}

	params := sqlc.GetCertFromIdParams{
		ID:     certUuid,
		UserID: userId,
	}

	dbCert, errCert := state.Queries.GetCertFromId(ctx, params)
	if errCert != nil {
		return Certificate{}, errCert
	}

	certView := dbCert.ToCertView()
	apiCert := MapCertificateViewDbToDomain(certView)
	apiCert.calculateExpiryDate()

	return apiCert, nil

}

func GetCertificateFromListOfIds(state *internal.ApiState, ctx context.Context, certIds []uuid.UUID, userId uuid.UUID) ([]Certificate, error) {

	certs := make([]Certificate, 0, len(certIds))
	for _, certId := range certIds {
		cert, errCert := GetCertificateFromId(state, ctx, certId.String(), userId)
		if errCert != nil {
			return nil, errCert
		}
		certs = append(certs, cert)
	}
	return certs, nil

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

func (cert *Certificate) calculateExpiryDate() {

	if !cert.ManualExpiry.IsZero() {
		cert.ExpiryDate = cert.ManualExpiry
	} else if cert.CertType.NormalValidityMonths != 0 {
		cert.ExpiryDate = getExpiryAfterValidity(cert.IssuedDate, int(cert.CertType.NormalValidityMonths))
	} else {
		cert.ExpiryDate = time.Time{}
	}

}

func getExpiryAfterValidity(issueDate time.Time, validityMonths int) time.Time {

	issueDate = time.Date(
		issueDate.Year(),
		issueDate.Month(),
		issueDate.Day(),
		0, 0, 0, 0,
		issueDate.Location(),
	)

	target := issueDate.AddDate(0, int(validityMonths), 0)
	issueDay := issueDate.Day()
	daysInTargetMonth := daysInMonth(target.Year(), target.Month())
	if issueDay > daysInTargetMonth {
		issueDay = daysInTargetMonth
	}

	targetSameDay := time.Date(
		target.Year(),
		target.Month(),
		issueDay,
		0, 0, 0, 0,
		target.Location(),
	)

	return targetSameDay.AddDate(0, 0, -1)

}

func daysInMonth(year int, month time.Month) int {
	// 1st of next month, minus one day
	firstOfNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastOfThis := firstOfNext.AddDate(0, 0, -1)
	return lastOfThis.Day()
}
