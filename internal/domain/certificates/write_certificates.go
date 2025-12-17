package certificates

import (
	"context"
	"errors"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/apiHttp/auth"
	"github.com/adamjames870/seacert/internal/database/sqlc"
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

	user, okUser := auth.UserFromContext(ctx)
	if !okUser {
		return Certificate{}, errors.New("unable to get user from context")
	}

	newCert := sqlc.CreateCertParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     user.Id,
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

	return apiCert, nil

}
