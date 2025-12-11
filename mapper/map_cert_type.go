package mapper

import (
	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
)

func MapCertificateType(certType database.CertificateType) models.CertificateType {

	return models.CertificateType{
		Id:                   certType.ID,
		CreatedAt:            certType.CreatedAt,
		UpdatedAt:            certType.UpdatedAt,
		Name:                 certType.Name,
		ShortName:            certType.ShortName,
		StcwReference:        certType.StcwReference.String,
		NormalValidityMonths: certType.NormalValidityMonths.Int32,
	}
}
