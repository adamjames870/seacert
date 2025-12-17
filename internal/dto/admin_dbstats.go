package dto

type DbStats struct {
	CountCert     int    `json:"count-certs"`
	CountCertType int    `json:"count-cert-types"`
	CountIssuer   int    `json:"count-issuers"`
	CountUsers    int    `json:"count-users"`
	UserId        string `json:"user-id"`
	UserEmail     string `json:"user-email"`
}
