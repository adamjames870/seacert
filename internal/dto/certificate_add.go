package dto

type ParamsAddCertificate struct {
	CertTypeId      string `json:"cert-type-id" validate:"required"`
	CertNumber      string `json:"cert-number" validate:"required"`
	Issuer          string `json:"issuer" validate:"required"`
	IssuedDate      string `json:"issued-date" validate:"required"`
	AlternativeName string `json:"alternative-name"`
	Remarks         string `json:"remarks"`
}
