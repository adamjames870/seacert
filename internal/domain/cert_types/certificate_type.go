package cert_types

import (
	"time"

	"github.com/google/uuid"
)

type SuccessionReason string

const (
	ReasonUpdated  SuccessionReason = "updated"
	ReasonReplaced SuccessionReason = "replaced"
)

type CertificateType struct {
	Id                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"created-at"`
	UpdatedAt            time.Time `json:"updated-at"`
	Name                 string    `json:"name"`
	ShortName            string    `json:"short_name"`
	StcwReference        string    `json:"stcw_reference"`
	NormalValidityMonths int32     `json:"normal_validity_months"`
}
