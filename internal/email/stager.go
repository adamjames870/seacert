package email

import (
	"github.com/adamjames870/seacert/internal/email/email_templates"
)

type Request struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

func Welcome(request Request) (sentId string, err error) {

	data := email_templates.WelcomeEmailData{
		FirstName: request.Name,
		AppURL:    "https://www.seacert.app/",
	}

	mail, errTmplt := email_templates.GetWelcomeEmail(data, []string{request.To})

	if errTmplt != nil {
		return "", errTmplt
	}

	id, errSend := Send(mail)
	if errSend != nil {
		return "", errSend
	}

	return id, nil

}
