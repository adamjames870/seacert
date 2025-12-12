package certificates

import (
	"time"

	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/google/uuid"
)

type Certificate struct {
	ID              uuid.UUID                  `json:"id"`
	CreatedAt       time.Time                  `json:"created-at"`
	UpdatedAt       time.Time                  `json:"updated-at"`
	CertType        cert_types.CertificateType `json:"cert-type-id"`
	CertNumber      string                     `json:"cert-number"`
	Issuer          string                     `json:"issuer"`
	IssuedDate      time.Time                  `json:"issued-date"`
	AlternativeName string                     `json:"alternative-name"`
	Remarks         string                     `json:"remarks"`
}
