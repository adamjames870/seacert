package dto

import "time"

type Certificate struct {
	Id                string    `json:"id"`
	CreatedAt         time.Time `json:"created-at"`
	UpdatedAt         time.Time `json:"updated-at"`
	CertTypeId        string    `json:"cert-type-id"`
	CertTypeName      string    `json:"cert-type-name"`
	CertTypeShortName string    `json:"cert-type-short-name"`
	CertTypeStcwRef   string    `json:"cert-type-stcw-ref"`
	CertNumber        string    `json:"cert-number"`
	IssuerId          string    `json:"issuer-id"`
	IssuerName        string    `json:"issuer-name"`
	IssuerCountry     string    `json:"issuer-country"`
	IssuerWebsite     string    `json:"issuer-website"`
	IssuedDate        time.Time `json:"issued-date"`
	ExpiryDate        time.Time `json:"expiry-date"`
	AlternativeName   string    `json:"alternative-name"`
	Remarks           string    `json:"remarks"`
}

type ParamsAddCertificate struct {
	UserId          string  `json:"user-id" validate:"required"`
	CertTypeId      string  `json:"cert-type-id" validate:"required"`
	CertNumber      string  `json:"cert-number" validate:"required"`
	IssuerId        string  `json:"issuer-id" validate:"required"`
	IssuedDate      string  `json:"issued-date" validate:"required"`
	AlternativeName *string `json:"alternative-name,omitempty"`
	Remarks         *string `json:"remarks,omitempty"`
	ManualExpiry    *string `json:"manual-expiry,omitempty"`
}

type ParamsUpdateCertificate struct {
	UserId          string  `json:"user-id"`
	Id              string  `json:"id" validate:"required"`
	CertNumber      *string `json:"cert-number,omitempty"`
	CertTypeId      *string `json:"cert-type-id,omitempty"`
	IssuerId        *string `json:"issuer-id,omitempty"`
	IssuedDate      *string `json:"issued-date,omitempty"`
	AlternativeName *string `json:"alternative-name,omitempty"`
	Remarks         *string `json:"remarks,omitempty"`
	ManualExpiry    *string `json:"manual-expiry,omitempty"`
}
