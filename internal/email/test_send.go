package email

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/email/templates"
	"github.com/resend/resend-go/v2"
)

const senderEmail = "SeaCert <notifications@seacert.app>"

type TestEmailRequest struct {
	To string `json:"to"`
}

func TestEmail(state *internal.ApiState, request TestEmailRequest, w http.ResponseWriter, r *http.Request) {

	data := templates.TestEmailData{
		Title:     "SeaCert Test Email",
		Year:      2026,
		FirstName: "Adam",
		AppURL:    "https://www.seacert.app",
	}

	email, err := templates.GetTestEmail(data, []string{request.To})
	if err != nil {
		http.Error(w, "Failed to get email template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      email.Recipient,
		Subject: email.Title,
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
