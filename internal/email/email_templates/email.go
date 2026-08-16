package email_templates

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
