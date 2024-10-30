package util

import (
	"Take_Out/config"
	"gopkg.in/mail.v2"
)

type EmailSender struct {
}

func SendEmail(receive, data, subject string) error {
	m := mail.NewMessage()
	m.SetHeader("From", config.Config.Email.SmtpEmail)
	m.SetHeader("To", receive)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)
	d := mail.NewDialer(config.Config.Email.SmtpHost, 465, config.Config.Email.SmtpEmail, config.Config.Email.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
