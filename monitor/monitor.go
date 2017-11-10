package monitor

import (
	"net/smtp"
	"strings"

	"github.com/domodwyer/mailyak"
)

func NewEmailProcessor(smtpUser string, smtpPass string, smtpServer string) EmailManipulator {

	mail := mailyak.New(smtpServer, smtp.PlainAuth(
		"",
		smtpUser,
		smtpPass,
		strings.Split(smtpServer, ":")[0],
	))

	return &Email{
		smtpUser:   smtpUser,
		smtpPass:   smtpPass,
		smtpServer: smtpServer,
		email:      mail,
	}
}

func (e *Email) SendEmail(from string, to string) error {
	e.email.To(to)
	e.email.From(to)
	e.email.Subject("")

	e.email.Plain().Set("")

	if err := e.email.Send(); err != nil {
		return err
	}
	return nil
}
