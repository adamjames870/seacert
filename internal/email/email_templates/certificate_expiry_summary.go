package email_templates

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"
)

type ExpiringCertificateItem struct {
	Name       string
	ExpiryDate string
	Url        string
}

type CertificateExpiryEmailData struct {
	FirstName     string
	AppURL        string
	Expired       []ExpiringCertificateItem
	NextMonth     []ExpiringCertificateItem
	NextSixMonths []ExpiringCertificateItem
	NextYear      []ExpiringCertificateItem
}

type certificateExpiryEmailTemplateData struct {
	BaseEmailData
	CertificateExpiryEmailData
}

func GetCertificateExpiryEmail(
	data CertificateExpiryEmailData,
	recipient []string,
) (Email, error) {

	templateData := certificateExpiryEmailTemplateData{
		CertificateExpiryEmailData: data,
		BaseEmailData: BaseEmailData{
			Title: "Your SeaCert certificate expiry summary",
			Year:  time.Now().Year(),
		},
	}

	tmpl, err := template.ParseFS(
		templateFS,
		"base.html",
		"certificate_expiry_summary.html",
	)
	if err != nil {
		return Email{}, err
	}

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base.html", templateData)
	if err != nil {
		return Email{}, err
	}

	textBody := buildCertificateExpiryText(data, templateData.Year)

	return Email{
		Subject:   "Your SeaCert certificate expiry summary",
		HtmlBody:  buf.String(),
		TextBody:  textBody,
		Recipient: recipient,
	}, nil
}

func buildCertificateExpiryText(
	data CertificateExpiryEmailData,
	year int,
) string {

	var b strings.Builder

	fmt.Fprintf(&b,
		"Certificate expiry summary\n\nHello %s,\n\n"+
			"Here is your current SeaCert certificate expiry summary.\n\n",
		data.FirstName,
	)

	writeCertificateGroup := func(
		heading string,
		certificates []ExpiringCertificateItem,
	) {
		if len(certificates) == 0 {
			return
		}

		fmt.Fprintf(&b, "%s\n", heading)

		for _, cert := range certificates {
			fmt.Fprintf(
				&b,
				"- %s — %s\n",
				cert.Name,
				cert.ExpiryDate,
			)
		}

		b.WriteString("\n")
	}

	writeCertificateGroup(
		"EXPIRED",
		data.Expired,
	)

	writeCertificateGroup(
		"EXPIRING WITHIN THE NEXT MONTH",
		data.NextMonth,
	)

	writeCertificateGroup(
		"EXPIRING WITHIN THE NEXT 6 MONTHS",
		data.NextSixMonths,
	)

	writeCertificateGroup(
		"EXPIRING WITHIN THE NEXT YEAR",
		data.NextYear,
	)

	fmt.Fprintf(
		&b,
		"View your certificates:\n%s\n\n© %d SeaCert\n",
		data.AppURL,
		year,
	)

	return b.String()
}
