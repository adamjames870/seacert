package dto

type cert_type_succession struct {
	Id              string `json:"id"`
	ReplacingType   string `json:"replacing-type-id"`
	ReplaceableType string `json:"replaceable-type-id"`
	ReplaceReason   string `json:"replace-reason"`
}

type cert_type_successions struct {
	Successions []cert_type_succession `json:"successions"`
}

type ParamsAddCertTypeSuccession struct {
	ReplacingType   string `json:"replacing-type-id"`
	ReplaceableType string `json:"replaceable-type-id"`
	ReplaceReason   string `json:"replace-reason"`
}
