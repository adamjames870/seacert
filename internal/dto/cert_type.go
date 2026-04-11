package dto

import "time"

type CertificateType struct {
	Id                   string    `json:"id"`
	CreatedAt            time.Time `json:"created-at"`
	UpdatedAt            time.Time `json:"updated-at"`
	Name                 string    `json:"name"`
	ShortName            string    `json:"short-name"`
	StcwRef              string    `json:"stcw-reference"`
	NormalValidityMonths int32     `json:"normal-validity-months"`
	Status               string    `json:"status"`
	CreatedBy            *string   `json:"created-by,omitempty"`
}

type ParamsAddCertificateType struct {
	Name                 string  `json:"name" validate:"required"`
	ShortName            string  `json:"short-name" validate:"required"`
	StcwReference        *string `json:"stcw-reference,omitempty"`
	NormalValidityMonths *int32  `json:"normal-validity-months" validate:"required"`
}

type ParamsUpdateCertificateType struct {
	Id                   string  `json:"id"`
	Name                 *string `json:"name,omitempty"`
	ShortName            *string `json:"short-name,omitempty"`
	StcwReference        *string `json:"stcw-reference,omitempty"`
	NormalValidityMonths *int32  `json:"normal-validity-months,omitempty"`
	Status               *string `json:"status,omitempty"`
}

type ParamsResolveCertificateType struct {
	ProvisionalId string `json:"provisional-id" validate:"required"`
	ReplacementId string `json:"replacement-id" validate:"required"`
}
