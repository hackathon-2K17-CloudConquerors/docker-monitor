package monitor

import (
	"github.com/domodwyer/mailyak"
)

type Email struct {
	smtpUser   string
	smtpPass   string
	smtpServer string
	email      *mailyak.MailYak
}
