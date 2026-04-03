package domain

const (
	// Roles
	RoleUser  = "user"
	RoleAdmin = "admin"

	// Certificate Type Status
	StatusProvisional = "provisional"
	StatusApproved    = "approved"

	// Succession Reasons (matching database enum)
	ReasonSupersededByCorrection = "superseded-by-correction"
	ReasonSupersededByRenewal    = "superseded-by-renewal"
)
