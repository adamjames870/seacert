package email_templates

import "embed"

//go:embed *.html
var templateFS embed.FS

const baseTemplate = "internal/email/email_templates/base.html"

type BaseEmailData struct {
	Title string
	Year  int
}

type Email struct {
	Subject   string
	HtmlBody  string
	TextBody  string
	Recipient []string
}
