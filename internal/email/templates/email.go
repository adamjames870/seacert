package templates

const baseTemplate = "internal/email/templates/base.html"

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
