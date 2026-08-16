package email

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/email/templates"
)

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

	sendEmail(email, w, r)

}
