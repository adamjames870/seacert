package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

const baseTemplate = "internal/email/templates/base.html"

type TestEmailData struct {
	Title     string
	Year      int
	FirstName string
	AppURL    string
}

type Email struct {
	Subject   string
	HtmlBody  string
	TextBody  string
	Recipient []string
}

func GetTestEmail(data TestEmailData, recipient []string) (Email, error) {

	tmpl, err := template.ParseFiles(baseTemplate, "internal/email/templates/test_email.html")
	if err != nil {
		return Email{}, err
	}

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return Email{}, err
	}

	htmlBody := buf.String()

	textBody := fmt.Sprintf(
		`SeaCert

Hello %s,

This is a test email sent from the SeaCert API.

If you've received this message, outbound email delivery is configured correctly.

Open SeaCert:
%s

© %d SeaCert
`,
		data.FirstName,
		data.AppURL,
		data.Year,
	)

	return Email{
		Subject:   data.Title,
		HtmlBody:  htmlBody,
		TextBody:  textBody,
		Recipient: recipient,
	}, nil

}
