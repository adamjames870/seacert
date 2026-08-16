package email

import (
	"bytes"
	"html/template"
)

type TemplateData struct {
	Title   string
	Heading string
	Message string
	CTA     string
	CTAURL  string
}

func renderHTML(data TemplateData) (string, error) {
	tmpl, err := template.ParseFiles(
		"internal/email/templates/layout.html",
	)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
