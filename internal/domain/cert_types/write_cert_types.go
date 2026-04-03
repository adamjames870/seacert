package cert_types

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func GetCertTypes(ctx context.Context, repo domain.Repository, userId *uuid.UUID, isAdmin bool) ([]CertificateType, error) {
	var certTypes []sqlc.CertificateType
	var err error

	if isAdmin {
		certTypes, err = repo.GetCertTypes(ctx)
	} else if userId != nil {
		certTypes, err = repo.GetCertTypesForUser(ctx, uuid.NullUUID{UUID: *userId, Valid: true})
	} else {
		certTypes, err = repo.GetCertTypesForUser(ctx, uuid.NullUUID{Valid: false})
	}

	if err != nil {
		return nil, err
	}

	apiCertTypes := make([]CertificateType, 0, len(certTypes))
	for _, cType := range certTypes {
		apiCertTypes = append(apiCertTypes, MapCertificateTypeDbToDomain(cType))
	}

	return apiCertTypes, nil
}

func GetCertTypeFromId(ctx context.Context, repo domain.Repository, id uuid.UUID) (CertificateType, error) {
	certType, err := repo.GetCertTypeFromId(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CertificateType{}, domain.ErrNotFound
		}
		return CertificateType{}, err
	}

	return MapCertificateTypeDbToDomain(certType), nil
}

func GetCertTypeFromName(ctx context.Context, repo domain.Repository, name string) (CertificateType, error) {
	certType, err := repo.GetCertTypeFromName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CertificateType{}, domain.ErrNotFound
		}
		return CertificateType{}, err
	}

	return MapCertificateTypeDbToDomain(certType), nil
}

func CreateCertType(ctx context.Context, repo domain.Repository, params dto.ParamsAddCertificateType, creatorId uuid.UUID, isAdmin bool) (CertificateType, error) {
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
		StcwReference:        domain.ToNullStringFromPointer(params.StcwReference),
		NormalValidityMonths: domain.ToNullInt32OrNil(params.NormalValidityMonths),
		Status:               status,
		CreatedBy: uuid.NullUUID{
			UUID:  creatorId,
			Valid: creatorId != uuid.Nil,
		},
	}

	dbCertType, err := repo.CreateCertType(ctx, newCert)
	if err != nil {
		return CertificateType{}, err
	}

	return MapCertificateTypeDbToDomain(dbCertType), nil
}

func UpdateCertificateType(ctx context.Context, repo domain.Repository, params dto.ParamsUpdateCertificateType) (CertificateType, error) {
	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return CertificateType{}, domain.ErrInvalidInput
	}

	updateCert := sqlc.UpdateCertTypeParams{
		ID:                   uuidId,
		Name:                 domain.ToNullStringFromPointer(params.Name),
		ShortName:            domain.ToNullStringFromPointer(params.ShortName),
		StcwReference:        domain.ToNullStringFromPointer(params.StcwReference),
		NormalValidityMonths: domain.ToNullInt32FromPointer(params.NormalValidityMonths),
		Status:               domain.ToNullStringFromPointer(params.Status),
	}

	dbCertType, err := repo.UpdateCertType(ctx, updateCert)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CertificateType{}, domain.ErrNotFound
		}
		return CertificateType{}, err
	}

	return MapCertificateTypeDbToDomain(dbCertType), nil
}

func ResolveProvisionalCertType(ctx context.Context, repo domain.Repository, params dto.ParamsResolveCertificateType) error {
	provisionalId, errProv := uuid.Parse(params.ProvisionalId)
	if errProv != nil {
		return domain.ErrInvalidInput
	}

	replacementId, errRepl := uuid.Parse(params.ReplacementId)
	if errRepl != nil {
		return domain.ErrInvalidInput
	}

	return repo.WithTx(ctx, func(txRepo domain.Repository) error {
		err := txRepo.UpdateCertTypeReferences(ctx, sqlc.UpdateCertTypeReferencesParams{
			CertTypeID:   provisionalId,
			CertTypeID_2: replacementId,
		})
		if err != nil {
			return err
		}

		return txRepo.DeleteCertType(ctx, provisionalId)
	})
}
