package cert_types

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewCertType(state *internal.ApiState, ctx context.Context, params dto.ParamsAddCertificateType) (CertificateType, error) {

	stcwRef := domain.ToNullStringFromPointer(params.StcwReference)
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

	dbCertType, errWriteNewCertType := state.Queries.CreateCertType(ctx, newCert)
	if errWriteNewCertType != nil {
		return CertificateType{}, errWriteNewCertType
	}

	apiCertType := MapCertificateTypeDbToDomain(dbCertType)

	return apiCertType, nil

}
