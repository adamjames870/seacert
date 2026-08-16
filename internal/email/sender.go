package email

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal/email/email_templates"
	"github.com/resend/resend-go/v2"
)

const senderEmail = "SeaCert <notifications@seacert.app>"

func sendEmail(email email_templates.Email, w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": sent.Id,
	})

}
