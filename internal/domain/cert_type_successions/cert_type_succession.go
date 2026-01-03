package cert_type_successions

import (
	"time"

	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/google/uuid"
)

type CertTypeSuccession struct {
	Id              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ReplacingType   cert_types.CertificateType
	ReplaceableType cert_types.CertificateType
	ReplaceReason   domain.SuccessionReason
}
