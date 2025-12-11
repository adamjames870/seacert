package db

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func WriteNewCertType(context context.Context, db database.Queries, params models.ParamsAddCertificateType) (database.CertificateType, error) {

	stcwRef := toNullString(params.StcwReference)
	normalValidity := toNullInt32OrNil(params.NormalValidityMonths)

	newCert := database.CreateCertTypeParams{
		ID:                   uuid.New(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Name:                 params.Name,
		ShortName:            params.ShortName,
		StcwReference:        stcwRef,
		NormalValidityMonths: normalValidity,
	}

	return db.CreateCertType(context, newCert)
}
