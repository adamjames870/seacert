package email_templates

import (
	"bytes"
	"fmt"
	"html/template"
	stdTime "time"
)

type WelcomeEmailData struct {
	FirstName string
	AppURL    string
}

type welcomeEmailData struct {
	BaseEmailData
	WelcomeEmailData
}

func GetWelcomeEmail(
	data WelcomeEmailData,
	recipient []string,
) (Email, error) {

	emailData := welcomeEmailData{
		BaseEmailData: BaseEmailData{
			Title: "Welcome to SeaCert",
			Year:  stdTime.Now().Year(),
		},
		WelcomeEmailData: data,
	}

	tmpl, err := template.ParseFS(
		templateFS,
		"base.html",
		"welcome.html",
	)
	if err != nil {
		return Email{}, err
	}

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base.html", emailData)
	if err != nil {
		return Email{}, err
	}

	textBody := fmt.Sprintf(
		`Welcome to SeaCert

Hello %s,

Welcome to SeaCert. Your account is ready to use.

SeaCert helps you keep track of your maritime certificates, monitor expiry dates, and keep your certification records organised in one place.

A good place to start is by adding your current certificates and their expiry dates.

Open SeaCert:
%s

© %d SeaCert
`,
		emailData.FirstName,
		emailData.AppURL,
		emailData.Year,
	)

	return Email{
		Subject:   "Welcome to SeaCert",
		HtmlBody:  buf.String(),
		TextBody:  textBody,
		Recipient: recipient,
	}, nil
}
