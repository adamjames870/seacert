package certificates

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
)

func MapCertificate(cert sqlc.Certificate, certType sqlc.CertificateType) Certificate {

	certModelType := cert_types.MapCertificateType(certType)

	return Certificate{
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
