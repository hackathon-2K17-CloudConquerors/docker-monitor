package monitor

import (
	"net/smtp"
	"strings"

	"github.com/domodwyer/mailyak"
)

func newEmailProcessor(smtpUser string, smtpPass string, smtpServer string, from string, to string) *email {

	mail := mailyak.New(smtpServer, smtp.PlainAuth(
		"",
		smtpUser,
		smtpPass,
		strings.Split(smtpServer, ":")[0],
	))

	return &email{
		smtpUser:   smtpUser,
		smtpPass:   smtpPass,
		smtpServer: smtpServer,
		from:       from,
		to:         to,
		email:      mail,
	}
}

func (e *email) sendEmail(content *emailcontent) error {
	e.Lock()
	e.email.To(e.to)
	e.email.From(e.from)
	e.email.Subject("ContainerEvents")

	e.email.Plain().Set("ContainerEvents" + content.ContainerID + "  " + content.ContainerName + "  " + content.time)

	if err := e.email.Send(); err != nil {
		return err
	}
	e.Unlock()

	return nil
}
