package models

import (
	"time"

	"github.com/google/uuid"
)

type ParamsAddCertificate struct {
	Name       string `json:"name" validate:"required"`
	CertNumber string `json:"cert-number" validate:"required"`
	Issuer     string `json:"issuer" validate:"required"`
	IssuedDate string `json:"issued-date" validate:"required"`
}

type Certificate struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created-at"`
	UpdatedAt  time.Time `json:"updated-at"`
	Name       string    `json:"name"`
	CertNumber string    `json:"cert-number"`
	Issuer     string    `json:"issuer"`
	IssuedDate time.Time `json:"issued-date"`
}
