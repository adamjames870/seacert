package certificates

import (
	"context"
	"errors"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewCert(state *internal.ApiState, ctx context.Context, params dto.ParamsAddCertificate) (Certificate, error) {

	issuedDate, errParse := time.Parse("2006-01-02", params.IssuedDate)
	if errParse != nil {
		return Certificate{}, errParse
	}

	certTypeId := params.CertTypeId
	apiCertType, errGetCertType := cert_types.GetCertTypeFromId(state, ctx, certTypeId)
	if errGetCertType != nil {
		return Certificate{}, errors.New("Error loading cert type: " + errGetCertType.Error())
	}

	issuerId := params.IssuerId
	apiIssuer, errGetIssuer := issuers.GetIssuerFromId(state, ctx, issuerId)
	if errGetIssuer != nil {
		return Certificate{}, errors.New("Error loading issuer: " + errGetIssuer.Error())
	}

	uuidId, errParse := uuid.Parse(params.UserId)
	if errParse != nil {
		return Certificate{}, errParse
	}

	timeNow := time.Now()

	newCert := sqlc.CreateCertParams{
		ID:         uuid.New(),
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		UserID:     uuidId,
		CertTypeID: apiCertType.Id,
		CertNumber: params.CertNumber,
		IssuerID:   apiIssuer.Id,
		IssuedDate: issuedDate,
	}

	dbCert, errCreateCert := state.Queries.CreateCert(ctx, newCert)
	if errCreateCert != nil {
		return Certificate{}, errCreateCert
	}

	apiCert := MapCertificateDbToDomain(dbCert, apiCertType, apiIssuer)

	if params.Supersedes != nil {

		uuidSupersedes, errParse := uuid.Parse(*params.Supersedes)
		if errParse != nil {
			return Certificate{}, errParse
		}

		if params.SupersedeReason == nil {
			return Certificate{}, errors.New("supersede reason is required")
		}

		predecessorCert, errPredecessor := GetCertificateFromId(state, ctx, uuidSupersedes.String(), uuidId)
		if errPredecessor != nil {
			return Certificate{}, errPredecessor
		}

		paramsSuccession := sqlc.CreateSuccessionParams{
			ID:      uuid.New(),
			OldCert: predecessorCert.Id,
			NewCert: dbCert.ID,
			Reason:  sqlc.SuccessionReason(*params.SupersedeReason),
		}

		_, errSuccession := state.Queries.CreateSuccession(ctx, paramsSuccession)
		if errSuccession != nil {
			return Certificate{}, errSuccession
		}

		predecessor := Predecesor{
			Cert:          predecessorCert,
			ReplaceReason: domain.SuccessionReason(paramsSuccession.Reason),
		}

		apiCert.Predecessors = append(apiCert.Predecessors, predecessor)

	}

	apiCert.calculateExpiryDate()

	return apiCert, nil

}

func UpdateCertificate(state *internal.ApiState, ctx context.Context, params dto.ParamsUpdateCertificate) (Certificate, error) {

	userId, errParse := uuid.Parse(params.UserId)
	if errParse != nil {
		return Certificate{}, errParse
	}

	_, errDb := GetCertificateFromId(state, ctx, params.Id, userId)
	if errDb != nil {
		return Certificate{}, errDb
	}

	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return Certificate{}, errParse
	}

	if params.CertTypeId != nil {
		_, errCertType := cert_types.GetCertTypeFromId(state, ctx, *params.CertTypeId)
		if errCertType != nil {
			return Certificate{}, errCertType
		}
	}

	if params.IssuerId != nil {
		_, errIssuer := issuers.GetIssuerFromId(state, ctx, *params.IssuerId)
		if errIssuer != nil {
			return Certificate{}, errIssuer
		}
	}

	updatedCertificate := sqlc.UpdateCertificateParams{
		ID:              uuidId,
		CertNumber:      domain.ToNullStringFromPointer(params.CertNumber),
		IssuedDate:      domain.ToNullTimeFromStringPointer(params.IssuedDate),
		CertTypeID:      domain.ToNullUUIDFromStringPointer(params.CertTypeId),
		AlternativeName: domain.ToNullStringFromPointer(params.AlternativeName),
		Remarks:         domain.ToNullStringFromPointer(params.Remarks),
		IssuerID:        domain.ToNullUUIDFromStringPointer(params.IssuerId),
		Deleted:         domain.ToNullBoolFromPointer(params.Deleted),
	}

	dbCert, errUpdateCert := state.Queries.UpdateCertificate(ctx, updatedCertificate)
	if errUpdateCert != nil {
		return Certificate{}, errUpdateCert
	}

	certType, errCertType := cert_types.GetCertTypeFromId(state, ctx, dbCert.CertTypeID.String())
	if errCertType != nil {
		return Certificate{}, errCertType
	}

	issuer, errIssuer := issuers.GetIssuerFromId(state, ctx, dbCert.IssuerID.String())
	if errIssuer != nil {
		return Certificate{}, errIssuer
	}

	apiCert := MapCertificateDbToDomain(dbCert, certType, issuer)
	apiCert.calculateExpiryDate()

	return apiCert, nil

}
