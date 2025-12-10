package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID              uuid.UUID       `json:"id"`
	CreatedAt       time.Time       `json:"created-at"`
	UpdatedAt       time.Time       `json:"updated-at"`
	CertType        CertificateType `json:"name"`
	CertNumber      string          `json:"cert-number"`
	Issuer          string          `json:"issuer"`
	IssuedDate      time.Time       `json:"issued-date"`
	AlternativeName string          `json:"alternative-name"`
	Remarks         string          `json:"remarks"`
}

type ParamsAddCertificate struct {
	CertType        uuid.UUID `json:"name"`
	CertNumber      string    `json:"cert-number"`
	Issuer          string    `json:"issuer"`
	IssuedDate      string    `json:"issued-date"`
	AlternativeName string    `json:"alternative-name"`
	Remarks         string    `json:"remarks"`
}
