package dto

import "time"

type Certificate struct {
	Id                           string    `json:"id"`
	CreatedAt                    time.Time `json:"created-at"`
	UpdatedAt                    time.Time `json:"updated-at"`
	CertTypeId                   string    `json:"cert-type-id"`
	CertTypeName                 string    `json:"cert-type-name"`
	CertTypeShortName            string    `json:"cert-type-short-name"`
	CertTypeStcwRef              string    `json:"cert-type-stcw-ref"`
	CertTypeNormalValidityMonths int32     `json:"cert-type-normal-validity-months"`
	CertNumber                   string    `json:"cert-number"`
	IssuerId                     string    `json:"issuer-id"`
	IssuerName                   string    `json:"issuer-name"`
	IssuerCountry                string    `json:"issuer-country"`
	IssuerWebsite                string    `json:"issuer-website"`
	IssuedDate                   time.Time `json:"issued-date"`
	AlternativeName              string    `json:"alternative-name"`
	Remarks                      string    `json:"remarks"`
}
