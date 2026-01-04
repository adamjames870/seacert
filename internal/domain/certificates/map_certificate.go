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
		Id:              cert.ID,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertType:        certType,
		CertNumber:      cert.CertNumber,
		Issuer:          issuer,
		IssuedDate:      cert.IssuedDate,
		AlternativeName: cert.AlternativeName.String,
		Remarks:         cert.Remarks.String,
		ManualExpiry:    cert.ManualExpiry.Time,
		Deleted:         cert.Deleted,
	}
}

func MapCertificateViewDbToDomain(dbCert sqlc.CertView) Certificate {

	certType := cert_types.CertificateType{
		Id:                   dbCert.CertTypeID,
		CreatedAt:            dbCert.CreatedAt,
		UpdatedAt:            dbCert.UpdatedAt,
		Name:                 dbCert.CertTypeName,
		ShortName:            dbCert.CertTypeShortName,
		StcwReference:        dbCert.CertTypeStcwReference.String,
		NormalValidityMonths: dbCert.NormalValidityMonths.Int32,
	}

	issuer := issuers.Issuer{
		Id:        dbCert.IssuerID,
		CreatedAt: dbCert.IssuerCreatedAt,
		UpdatedAt: dbCert.IssuerUpdatedAt,
		Name:      dbCert.IssuerName,
		Country:   dbCert.IssuerCountry.String,
		Website:   dbCert.IssuerWebsite.String,
	}

	apiCert := Certificate{
		Id:              dbCert.ID,
		CreatedAt:       dbCert.CreatedAt,
		UpdatedAt:       dbCert.UpdatedAt,
		CertType:        certType,
		CertNumber:      dbCert.CertNumber,
		Issuer:          issuer,
		IssuedDate:      dbCert.IssuedDate,
		AlternativeName: dbCert.AlternativeName.String,
		Remarks:         dbCert.Remarks.String,
		ManualExpiry:    dbCert.ManualExpiry.Time,
		Deleted:         dbCert.Deleted,
		HasSuccessors:   dbCert.HasSuccessor,
	}

	apiCert.calculateExpiryDate()
	return apiCert

}

func MapCertificateDomainToDto(cert Certificate) dto.Certificate {

	rv := dto.Certificate{
		Id:                cert.Id.String(),
		CreatedAt:         cert.CreatedAt,
		UpdatedAt:         cert.UpdatedAt,
		CertTypeId:        cert.CertType.Id.String(),
		CertTypeName:      cert.CertType.Name,
		CertTypeShortName: cert.CertType.ShortName,
		CertTypeStcwRef:   cert.CertType.StcwReference,
		CertNumber:        cert.CertNumber,
		IssuerId:          cert.Issuer.Id.String(),
		IssuerName:        cert.Issuer.Name,
		IssuerCountry:     cert.Issuer.Country,
		IssuerWebsite:     cert.Issuer.Website,
		IssuedDate:        cert.IssuedDate,
		ExpiryDate:        cert.ExpiryDate,
		AlternativeName:   cert.AlternativeName,
		Remarks:           cert.Remarks,
		Deleted:           cert.Deleted,
		Predecessors:      []dto.Predecessor{},
		HasSuccessors:     cert.HasSuccessors,
	}

	for _, predecessor := range cert.Predecessors {
		rv.Predecessors = append(rv.Predecessors, MapPredecessorDomainToDto(predecessor))
	}

	return rv

}

func MapCertificateDtoToDomain(cert dto.Certificate) Certificate {

	id, _ := uuid.Parse(cert.Id)

	certTypeDto := dto.CertificateType{
		Id:        cert.CertTypeId,
		Name:      cert.CertTypeName,
		ShortName: cert.CertTypeShortName,
		StcwRef:   cert.CertTypeStcwRef,
	}

	issuerDto := dto.Issuer{
		Name:    cert.IssuerName,
		Country: cert.IssuerCountry,
		Website: cert.IssuerWebsite,
	}

	certType := cert_types.MapCertificateTypeDtoToDomain(certTypeDto)
	issuer := issuers.MapIssuerDtoToDomain(issuerDto)

	rv := Certificate{
		Id:              id,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertType:        certType,
		CertNumber:      cert.CertNumber,
		Issuer:          issuer,
		IssuedDate:      cert.IssuedDate,
		AlternativeName: cert.AlternativeName,
		Remarks:         cert.Remarks,
		Deleted:         cert.Deleted,
		Predecessors:    []Predecesor{},
	}

	for _, predecessor := range cert.Predecessors {
		rv.Predecessors = append(rv.Predecessors, MapPredecessorDtoToDomain(predecessor))
	}

	return rv

}

func MapCertificateDomainToDb(cert Certificate) sqlc.Certificate {

	alternativeName := domain.ToNullString(cert.AlternativeName)
	remarks := domain.ToNullString(cert.Remarks)

	return sqlc.Certificate{
		ID:              cert.Id,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertNumber:      cert.CertNumber,
		IssuedDate:      cert.IssuedDate,
		CertTypeID:      cert.CertType.Id,
		AlternativeName: alternativeName,
		Remarks:         remarks,
		IssuerID:        cert.Issuer.Id,
		Deleted:         cert.Deleted,
	}

}

func MapPredecessorDomainToDto(predecessor Predecesor) dto.Predecessor {

	return dto.Predecessor{
		Cert:   MapCertificateDomainToDto(predecessor.Cert),
		Reason: string(predecessor.ReplaceReason),
	}

}

func MapPredecessorDtoToDomain(predecessor dto.Predecessor) Predecesor {
	return Predecesor{
		Cert:          MapCertificateDtoToDomain(predecessor.Cert),
		ReplaceReason: domain.SuccessionReason(predecessor.Reason),
	}
}
