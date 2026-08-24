package notifications

type Type string

const (
	TypeNoCertificates7Day       Type = "no_certificates_7d"
	TypeNoCertificates1Month     Type = "no_certificates_1m"
	TypeCertificateExpirySummary Type = "certificate_expiry_summary"
)
