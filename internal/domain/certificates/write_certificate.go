package certificates

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewCert(context context.Context, db sqlc.Queries, params dto.ParamsAddCertificate) (sqlc.Certificate, error) {

	issuedDate, errParse := time.Parse("2006-01-02", params.IssuedDate)
	if errParse != nil {
		return sqlc.Certificate{}, errParse
	}

	certTypeId, errParseCertTypeId := uuid.Parse(params.CertTypeId)
	if errParseCertTypeId != nil {
		return sqlc.Certificate{}, errParseCertTypeId
	}

	newCert := sqlc.CreateCertParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CertTypeID: certTypeId,
		CertNumber: params.CertNumber,
		Issuer:     params.Issuer,
		IssuedDate: issuedDate,
	}

	return db.CreateCert(context, newCert)
}
