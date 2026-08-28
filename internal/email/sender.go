package email

import (
	"os"

	"github.com/adamjames870/seacert/internal/email/email_templates"
	"github.com/resend/resend-go/v2"
)

const senderEmail = "SeaCert <notifications@seacert.app>"

func Send(email email_templates.Email) (sentId string, err error) {

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      email.Recipient,
		Subject: email.Subject,
		Html:    email.HtmlBody,
		Text:    email.TextBody,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return sent.Id, nil

}
