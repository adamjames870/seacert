package dto

import "time"

type User struct {
	Id                    string     `json:"id"`
	CreatedAt             time.Time  `json:"created-at"`
	UpdatedAt             time.Time  `json:"updated-at"`
	Forename              string     `json:"forename"`
	Surname               string     `json:"surname"`
	Email                 string     `json:"email"`
	Nationality           string     `json:"nationality"`
	Role                  string     `json:"role"`
	EmailConsent          bool       `json:"email_consent"`
	EmailConsentTimestamp *time.Time `json:"email_consent_timestamp"`
	EmailConsentVersion   string     `json:"email_consent_version"`
	EmailConsentSource    string     `json:"email_consent_source"`
}

type ParamsAddUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type ParamsUpdateUser struct {
	Id                  string  `json:"id"`
	Forename            *string `json:"forename,omitempty"`
	Surname             *string `json:"surname,omitempty"`
	Nationality         *string `json:"nationality,omitempty"`
	EmailConsent        *bool   `json:"email_consent,omitempty"`
	EmailConsentVersion *string `json:"email_consent_version,omitempty"`
	EmailConsentSource  *string `json:"email_consent_source,omitempty"`
}
