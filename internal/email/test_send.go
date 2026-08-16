package email

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/email/templates"
)

type EmailRequest struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

func TestEmail(state *internal.ApiState, request EmailRequest, w http.ResponseWriter, r *http.Request) {

	data := templates.NoCertsEmailData{
		FirstName: request.Name,
		AppURL:    "https://www.seacert.app",
	}

	email, err := templates.GetNoCertsEmail(data, []string{request.To})
	if err != nil {
		http.Error(w, "Failed to get email template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendEmail(email, w, r)

}
