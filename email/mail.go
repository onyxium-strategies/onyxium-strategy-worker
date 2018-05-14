package email

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/smtp"
)

func EmailActivateUser(to, id, token string) error {
	t, err := template.New("mail").Parse(tpl)
	if err != nil {
		log.Info("Parse ", err)
		return err
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
		return err
	}
	body := buf.String()
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Activate your account.\n"
	// to = "To: " + to + "\n"
	msg := []byte(subject + mime + "\n" + body)
	// mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// subject := "Activate your account.\n"
	// msg := []byte("To: " + to + "\r\n" +
	// 	"Subject: " + subject + "\n" +
	// 	mime + "\r\n" +
	// 	body + "\r\n")
	addr := "smtp.gmail.com:587"
	recipient := []string{to}
	auth := smtp.PlainAuth("", "alainfh94@gmail.com", "gfqodfvhlfwdmvft", "smtp.gmail.com")
	from := "alainfh94@gmail.com"
	err = smtp.SendMail(addr, auth, from, recipient, msg)
	return err
}
