package certificates

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func MapCertificateDbToDomain(cert sqlc.Certificate, certType cert_types.CertificateType, issuer issuers.Issuer) Certificate {

	return Certificate{
		ID:              cert.ID,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertType:        certType,
		CertNumber:      cert.CertNumber,
		Issuer:          issuer,
		IssuedDate:      cert.IssuedDate,
		AlternativeName: cert.AlternativeName.String,
		Remarks:         cert.Remarks.String,
	}
}

func MapCertificateDomainToDto(cert Certificate) dto.Certificate {

	return dto.Certificate{
		Id:                           cert.ID.String(),
		CreatedAt:                    cert.CreatedAt,
		UpdatedAt:                    cert.UpdatedAt,
		CertTypeId:                   cert.CertType.Id.String(),
		CertTypeName:                 cert.CertType.Name,
		CertTypeShortName:            cert.CertType.ShortName,
		CertTypeStcwRef:              cert.CertType.StcwReference,
		CertTypeNormalValidityMonths: cert.CertType.NormalValidityMonths,
		CertNumber:                   cert.CertNumber,
		IssuerId:                     cert.Issuer.Id.String(),
		IssuerName:                   cert.Issuer.Name,
		IssuerCountry:                cert.Issuer.Country,
		IssuerWebsite:                cert.Issuer.Website,
		IssuedDate:                   cert.IssuedDate,
		AlternativeName:              cert.AlternativeName,
		Remarks:                      cert.Remarks,
	}

}

func MapCertificateDtoToDomain(cert dto.Certificate) Certificate {

	id, _ := uuid.Parse(cert.Id)

	certTypeDto := dto.CertificateType{
		Id:                   cert.CertTypeId,
		Name:                 cert.CertTypeName,
		ShortName:            cert.CertTypeShortName,
		StcwRef:              cert.CertTypeStcwRef,
		NormalValidityMonths: cert.CertTypeNormalValidityMonths,
	}

	issuerDto := dto.Issuer{
		Name:    cert.IssuerName,
		Country: cert.IssuerCountry,
		Website: cert.IssuerWebsite,
	}

	certType := cert_types.MapCertificateTypeDtoToDomain(certTypeDto)
	issuer := issuers.MapIssuerDtoToDomain(issuerDto)

	return Certificate{
		ID:              id,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertType:        certType,
		CertNumber:      cert.CertNumber,
		Issuer:          issuer,
		IssuedDate:      cert.IssuedDate,
		AlternativeName: cert.AlternativeName,
		Remarks:         cert.Remarks,
	}

}

func MapCertificateDomainToDb(cert Certificate) sqlc.Certificate {

	alternativeName := domain.ToNullString(cert.AlternativeName)
	remarks := domain.ToNullString(cert.Remarks)

	return sqlc.Certificate{
		ID:              cert.ID,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertNumber:      cert.CertNumber,
		IssuerID:        cert.Issuer.Id,
		IssuedDate:      cert.IssuedDate,
		CertTypeID:      cert.CertType.Id,
		AlternativeName: alternativeName,
		Remarks:         remarks,
	}

}
