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

func WriteNewCertType(state *internal.ApiState, ctx context.Context, params dto.ParamsAddCertificateType, creatorId uuid.UUID, isAdmin bool) (CertificateType, error) {

	stcwRef := domain.ToNullStringFromPointer(params.StcwReference)
	normalValidity := domain.ToNullInt32OrNil(params.NormalValidityMonths)

	status := "provisional"
	if isAdmin {
		status = "approved"
	}

	newCert := sqlc.CreateCertTypeParams{
		ID:                   uuid.New(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Name:                 params.Name,
		ShortName:            params.ShortName,
		StcwReference:        stcwRef,
		NormalValidityMonths: normalValidity,
		Status:               status,
		CreatedBy: uuid.NullUUID{
			UUID:  creatorId,
			Valid: creatorId != uuid.Nil,
		},
	}

	dbCertType, errWriteNewCertType := state.Queries.CreateCertType(ctx, newCert)
	if errWriteNewCertType != nil {
		return CertificateType{}, errWriteNewCertType
	}

	apiCertType := MapCertificateTypeDbToDomain(dbCertType)

	return apiCertType, nil

}

func UpdateCertificateType(state *internal.ApiState, ctx context.Context, params dto.ParamsUpdateCertificateType) (CertificateType, error) {

	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return CertificateType{}, errParse
	}

	name := domain.ToNullStringFromPointer(params.Name)
	shortName := domain.ToNullStringFromPointer(params.ShortName)
	stcwRef := domain.ToNullStringFromPointer(params.StcwReference)
	normalValidity := domain.ToNullInt32FromPointer(params.NormalValidityMonths)
	status := domain.ToNullStringFromPointer(params.Status)

	updateCert := sqlc.UpdateCertTypeParams{
		ID:                   uuidId,
		Name:                 name,
		ShortName:            shortName,
		StcwReference:        stcwRef,
		NormalValidityMonths: normalValidity,
		Status:               status,
	}

	dbCertType, errUpdateCertType := state.Queries.UpdateCertType(ctx, updateCert)
	if errUpdateCertType != nil {
		return CertificateType{}, errUpdateCertType
	}

	apiCertType := MapCertificateTypeDbToDomain(dbCertType)
	return apiCertType, nil

}

func ResolveProvisionalCertType(state *internal.ApiState, ctx context.Context, params dto.ParamsResolveCertificateType) error {

	provisionalId, errProv := uuid.Parse(params.ProvisionalId)
	if errProv != nil {
		return errProv
	}

	replacementId, errRepl := uuid.Parse(params.ReplacementId)
	if errRepl != nil {
		return errRepl
	}

	// 1. Update all certificates using provisionalId to replacementId
	errUpdate := state.Queries.UpdateCertTypeReferences(ctx, sqlc.UpdateCertTypeReferencesParams{
		CertTypeID:   provisionalId,
		CertTypeID_2: replacementId,
	})
	if errUpdate != nil {
		return errUpdate
	}

	// 2. Delete the provisional cert type
	return state.Queries.DeleteCertType(ctx, provisionalId)
}
