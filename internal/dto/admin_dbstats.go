package dto

type DbStats struct {
	CountCert     int    `json:"count-certs"`
	CountCertType int    `json:"count-cert-types"`
	CountIssuer   int    `json:"count-issuers"`
	CountUsers    int    `json:"count-users"`
	UserId        string `json:"user-id"`
	UserEmail     string `json:"user-email"`
}

type NotificationsStatusEach struct {
	Type  string `json:"status-type"`
	Count int    `json:"count"`
}

type NotificationsStatusSummary struct {
	Statuses []NotificationsStatusEach `json:"statuses"`
}
