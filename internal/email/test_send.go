package email

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/resend/resend-go/v2"
)

type TestEmailRequest struct {
	To string `json:"to"`
}

func TestEmail(state *internal.ApiState, request TestEmailRequest, w http.ResponseWriter, r *http.Request) {

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From:    "SeaCert <notifications@seacert.app>",
		To:      []string{request.To},
		Subject: "SeaCert test email",
		Html: `
			<h1>Hello from SeaCert</h1>
			<p>If you're reading this, outbound email is working.</p>
		`,
		Text: "Hello from SeaCert. If you're reading this, outbound email is working.",
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
