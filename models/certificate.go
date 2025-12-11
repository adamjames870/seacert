package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID              uuid.UUID       `json:"id"`
	CreatedAt       time.Time       `json:"created-at"`
	UpdatedAt       time.Time       `json:"updated-at"`
	CertType        CertificateType `json:"cert-type-id"`
	CertNumber      string          `json:"cert-number"`
	Issuer          string          `json:"issuer"`
	IssuedDate      time.Time       `json:"issued-date"`
	AlternativeName string          `json:"alternative-name"`
	Remarks         string          `json:"remarks"`
}

type ParamsAddCertificate struct {
	CertTypeId      string `json:"cert-type-id" validate:"required"`
	CertNumber      string `json:"cert-number" validate:"required"`
	Issuer          string `json:"issuer" validate:"required"`
	IssuedDate      string `json:"issued-date" validate:"required"`
	AlternativeName string `json:"alternative-name"`
	Remarks         string `json:"remarks"`
}
