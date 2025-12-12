package dto

type ParamsAddCertificateType struct {
	Name                 string `json:"name"`
	ShortName            string `json:"short_name"`
	StcwReference        string `json:"stcw_reference"`
	NormalValidityMonths int32  `json:"normal_validity_months"`
}
