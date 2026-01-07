package dto

type Cert_type_successions struct {
	CertTypeId    string       `json:"cert-type-id"`
	CanReplace    []Succession `json:"can-replace-certs"`
	ReplaceableBy []Succession `json:"replaceable-by-certs"`
}

type Succession struct {
	CertType      CertificateType `json:"cert-type"`
	ReplaceReason string          `json:"replace-reason"`
}

type ParamsAddCertTypeSuccession struct {
	ReplacingType   string `json:"replacing-type-id"`
	ReplaceableType string `json:"replaceable-type-id"`
	ReplaceReason   string `json:"replace-reason"`
}
