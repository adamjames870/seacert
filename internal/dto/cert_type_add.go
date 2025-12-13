package dto

type ParamsAddCertificateType struct {
	Name                 string `json:"name"`
	ShortName            string `json:"short-name"`
	StcwReference        string `json:"stcw-reference"`
	NormalValidityMonths int32  `json:"normal-validity-months"`
}
