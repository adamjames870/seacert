package dto

type DbStats struct {
	CountCert     int    `json:"count-certs"`
	CountCertType int    `json:"count-cert-types"`
	CountIssuer   int    `json:"count-issuers"`
	UserId        string `json:"user-id"`
	UserEmail     string `json:"user-email"`
	UserRole      string `json:"user-role"`
}
