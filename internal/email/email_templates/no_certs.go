package email_templates

import (
	"bytes"
	"fmt"
	"html/template"
	stdTime "time"
)

type NoCertsEmailData struct {
	FirstName string
	AppURL    string
}

type noCertsEmailData struct {
	BaseEmailData
	NoCertsEmailData
}

func GetNoCertsEmail(
	data NoCertsEmailData,
	recipient []string,
) (Email, error) {

	emailData := noCertsEmailData{
		BaseEmailData: BaseEmailData{
			Title: "Get started with SeaCert",
			Year:  stdTime.Now().Year(),
		},
		NoCertsEmailData: data,
	}

	tmpl, err := template.ParseFiles(
		baseTemplate,
		"internal/email/email_templates/no_certs.html",
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
		`Get started with SeaCert

Hello %s,

We noticed that you haven't added any certificates to SeaCert yet.

Adding your certificates allows SeaCert to keep your certification records organised and help you stay ahead of upcoming expiry dates.

It only takes a few moments to add your first certificate.

Add your first certificate:
%s

Once your certificates are added, SeaCert can help you keep track of what you hold and when each certificate needs attention.

© %d SeaCert
`,
		emailData.FirstName,
		emailData.AppURL,
		emailData.Year,
	)

	return Email{
		Subject:   "Get started with SeaCert",
		HtmlBody:  buf.String(),
		TextBody:  textBody,
		Recipient: recipient,
	}, nil

}
