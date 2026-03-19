package cert_types

import (
	"errors"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func MapCertificateTypeDbToDomain(certType sqlc.CertificateType) CertificateType {

	return CertificateType{
		Id:                   certType.ID,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        certType.StcwReference.String,
		NormalValidityMonths: certType.NormalValidityMonths.Int32,
		Status:               certType.Status,
		CreatedBy:            certType.CreatedBy.UUID,
	}
}

func MapCertificateTypeDomainToDto(certType CertificateType) dto.CertificateType {

	var createdBy *string
	if certType.CreatedBy != uuid.Nil {
		s := certType.CreatedBy.String()
		createdBy = &s
	}

	return dto.CertificateType{
		Id:                   certType.Id.String(),
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwRef:              certType.StcwReference,
		NormalValidityMonths: certType.NormalValidityMonths,
		Status:               certType.Status,
		CreatedBy:            createdBy,
	}

}

func MapCertificateTypeDtoToDomain(certType dto.CertificateType) CertificateType {

	id, _ := uuid.Parse(certType.Id)
	var createdBy uuid.UUID
	if certType.CreatedBy != nil {
		createdBy, _ = uuid.Parse(*certType.CreatedBy)
	}

	return CertificateType{
		Id:                   id,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        certType.StcwRef,
		NormalValidityMonths: certType.NormalValidityMonths,
		Status:               certType.Status,
		CreatedBy:            createdBy,
	}

}

func MapCertificateTypeDomainToDb(certType CertificateType) sqlc.CertificateType {

	stcwRef := domain.ToNullString(certType.StcwReference)
	normalValidity := domain.ToNullInt32OrNil(certType.NormalValidityMonths)
	createdBy := uuid.NullUUID{
		UUID:  certType.CreatedBy,
		Valid: certType.CreatedBy != uuid.Nil,
	}

	return sqlc.CertificateType{
		ID:                   certType.Id,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        stcwRef,
		NormalValidityMonths: normalValidity,
		Status:               certType.Status,
		CreatedBy:            createdBy,
	}

}
func SuccessionReasonDbToDomain(reason sqlc.SuccessionReason) (SuccessionReason, error) {
	switch reason {
	case sqlc.SuccessionReplaced:
		return ReasonReplaced, nil
	case sqlc.SuccessionUpdated:
		return ReasonUpdated, nil
	}
	return "", errors.New("unknown succession reason")
}
