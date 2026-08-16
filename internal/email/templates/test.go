package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

import stdTime "time"

type TestEmailData struct {
	FirstName string
	AppURL    string
}

type testEmailData struct {
	BaseEmailData
	TestEmailData
}

func GetTestEmail(data TestEmailData, recipient []string) (Email, error) {

	emailData := testEmailData{
		BaseEmailData: BaseEmailData{
			Title: "SeaCert Test Email",
			Year:  stdTime.Now().Year(),
		},
		TestEmailData: data,
	}

	tmpl, err := template.ParseFiles(baseTemplate, "internal/email/templates/test.html")
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
		emailData.FirstName,
		emailData.AppURL,
		emailData.Year,
	)

	return Email{
		Subject:   emailData.Title,
		HtmlBody:  htmlBody,
		TextBody:  textBody,
		Recipient: recipient,
	}, nil

}
