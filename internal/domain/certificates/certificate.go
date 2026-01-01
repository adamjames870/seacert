package certificates

import (
	"time"

	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/google/uuid"
)

type Certificate struct {
	Id              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CertType        cert_types.CertificateType
	CertNumber      string
	Issuer          issuers.Issuer
	IssuedDate      time.Time
	ExpiryDate      time.Time
	AlternativeName string
	Remarks         string
	ManualExpiry    time.Time
	Deleted         bool
	Predecessors    []Predecesor
}

type Predecesor struct {
	Cert          Certificate
	ReplaceReason cert_types.SuccessionReason
}
