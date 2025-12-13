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
}
