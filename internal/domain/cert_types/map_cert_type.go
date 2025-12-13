package cert_types

import (
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
	}
}

func MapCertificateTypeDomainToDto(certType CertificateType) dto.CertificateType {

	return dto.CertificateType{
		Id:                   certType.Id.String(),
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwRef:              certType.StcwReference,
		NormalValidityMonths: certType.NormalValidityMonths,
	}

}

func MapCertificateTypeDtoToDomain(certType dto.CertificateType) CertificateType {

	id, _ := uuid.Parse(certType.Id)

	return CertificateType{
		Id:                   id,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        certType.StcwRef,
		NormalValidityMonths: certType.NormalValidityMonths,
	}

}

func MapCertificateTypeDomainToDb(certType CertificateType) sqlc.CertificateType {

	stcwRef := domain.ToNullString(certType.StcwReference)
	normalValidity := domain.ToNullInt32OrNil(certType.NormalValidityMonths)

	return sqlc.CertificateType{
		ID:                   certType.Id,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        stcwRef,
		NormalValidityMonths: normalValidity,
	}

}
