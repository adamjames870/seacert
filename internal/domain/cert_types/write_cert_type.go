package cert_types

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewCertType(context context.Context, db sqlc.Queries, params dto.ParamsAddCertificateType) (sqlc.CertificateType, error) {

	stcwRef := domain.ToNullString(params.StcwReference)
	normalValidity := domain.ToNullInt32OrNil(params.NormalValidityMonths)

	newCert := sqlc.CreateCertTypeParams{
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
