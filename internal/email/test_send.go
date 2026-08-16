package email

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/email/email_templates"
)

type EmailRequest struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

func TestEmail(state *internal.ApiState, request EmailRequest, w http.ResponseWriter, r *http.Request) {

	data := email_templates.CertificateExpiryEmailData{
		FirstName: request.Name,
		AppURL:    "https://www.seacert.app/certificates",

		Expired: []email_templates.ExpiringCertificateItem{
			{
				Name:       "ENG1 Medical Certificate",
				ExpiryDate: "02 Aug 2026",
				Url:        "https://www.seacert.app/certificates",
			},
		},

		NextMonth: []email_templates.ExpiringCertificateItem{
			{
				Name:       "GMDSS GOC",
				ExpiryDate: "08 Sep 2026",
				Url:        "https://www.seacert.app/certificates",
			},
			{
				Name:       "Basic STCW",
				ExpiryDate: "01 Jan 2021",
				Url:        "https://www.seacert.app/certificates",
			},
		},

		NextSixMonths: []email_templates.ExpiringCertificateItem{
			{
				Name:       "Advanced Fire Fighting",
				ExpiryDate: "14 Dec 2026",
				Url:        "https://www.seacert.app/certificates",
			},
		},

		NextYear: []email_templates.ExpiringCertificateItem{
			{
				Name:       "Master Unlimited",
				ExpiryDate: "12 May 2027",
				Url:        "https://www.seacert.app/certificates",
			},
			{
				Name:       "GMDSS",
				ExpiryDate: "12 May 2027",
				Url:        "https://www.seacert.app/certificates",
			},
		},
	}

	mail, err := email_templates.GetCertificateExpiryEmail(data, []string{request.To})

	if err != nil {
		http.Error(w, "Failed to get email template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendEmail(mail, w, r)

}
