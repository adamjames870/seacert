package cert_types

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
)

func MapCertificateType(certType sqlc.CertificateType) CertificateType {

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
