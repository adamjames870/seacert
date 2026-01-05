package cert_type_successions

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
)

func MapSuccessionDbToDomain(succession sqlc.CertificateTypeSuccession, replacing sqlc.CertificateType, replaceable sqlc.CertificateType) CertTypeSuccession {

	reason := 

	return CertTypeSuccession{
		Id:              succession.ID,
		CreatedAt:       succession.CreatedAt,
		UpdatedAt:       succession.UpdatedAt,
		ReplacingType:   cert_types.MapCertificateTypeDbToDomain(replacing),
		ReplaceableType: cert_types.MapCertificateTypeDbToDomain(replaceable),
		ReplaceReason:   succession.ReplaceReason,
	}

}
