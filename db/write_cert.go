package db

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func WriteNewCert(context context.Context, db database.Queries, params models.ParamsAddCertificate) (database.Certificate, error) {

	issuedDate, errParse := time.Parse("2006-01-02", params.IssuedDate)
	if errParse != nil {
		return database.Certificate{}, errParse
	}

	certTypeId, errParseCertTypeId := uuid.Parse(params.CertTypeId)
	if errParseCertTypeId != nil {
		return database.Certificate{}, errParseCertTypeId
	}

	newCert := database.CreateCertParams{
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
