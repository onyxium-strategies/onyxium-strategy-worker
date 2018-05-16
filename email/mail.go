package email

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
	"html/template"
)

const sendgridAPIKey = "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"

func EmailActivateUser(userEmail, id, token string) error {
	body, err := parseMailTemplate(id, token)
	if err != nil {
		return err
	}
	from := mail.NewEmail("Onyxium", "info@onyxium.io")
	subject := "Activate your account"
	to := mail.NewEmail("New user", userEmail)
	message := mail.NewSingleEmail(from, subject, to, body, body)
	client := sendgrid.NewSendClient(sendgridAPIKey)
	response, err := client.Send(message)

	// TODO needs to be converted to proper error handling
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return err
}

func parseMailTemplate(id, token string) (string, error) {
	t, err := template.New("mail").Parse(tpl)
	if err != nil {
		log.Info("Parse ", err)
		return "", err
	}
	templateData := struct {
		URL   string
		Id    string
		Token string
	}{
		URL:   "localhost:8000",
		Id:    id,
		Token: token,
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		log.Info("Execute ", err)
		return "", err
	}
	return buf.String(), nil
}
