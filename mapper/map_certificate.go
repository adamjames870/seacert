package mapper

import (
	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
)

func MapCertificate(cert database.Certificate, certType database.CertificateType) models.Certificate {

	certModelType := MapCertificateType(certType)

	return models.Certificate{
		ID:              cert.ID,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		CertType:        certModelType,
		CertNumber:      cert.CertNumber,
		Issuer:          cert.Issuer,
		IssuedDate:      cert.IssuedDate,
		AlternativeName: cert.AlternativeName.String,
		Remarks:         cert.Remarks.String,
	}
}
