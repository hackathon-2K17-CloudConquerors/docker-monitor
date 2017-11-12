package monitor

import (
	"net/smtp"
	"strconv"
	"strings"

	"github.com/domodwyer/mailyak"
	"go.uber.org/zap"
)

func newEmailProcessor(smtpUser string, smtpPass string, smtpServer string, to string) *email {

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
		to:         to,
		email:      mail,
	}
}

func (e *email) sendEmail(content *emailcontent) error {
	zap.L().Info("Sending Email", zap.Any("email content", content.container))
	e.Lock()
	e.email.To(e.to)
	e.email.From(e.smtpUser)
	e.email.Subject("Attention: " + content.container.ContainerName + " container down")

	time := strconv.Itoa(content.time.Hour()) + ":" + strconv.Itoa(content.time.Minute()) + ":" + strconv.Itoa(content.time.Second())
	day := content.time.Month().String() + " " + strconv.Itoa(content.time.Day()) + ", " + strconv.Itoa(content.time.Year())

	e.email.Plain().Set(
		"Application " + content.container.ContainerName + " went down at " + time + " on " + day + " (UTC). Please investigate\n\n" +
			"DETAILS: \n " +
			"ContainerName: " + content.container.ContainerName + "\n " +
			"ContainerID: " + content.container.ContainerID + " \n " +
			"ContainerNetwork: " + content.container.ContainerNetwork + " \n " +
			"Image: " + content.container.ImageName + " \n " +
			"CurrentStatus: " + content.container.Status + " \n " +
			"StartedAt: " + content.container.ContainerStatus + " \n " +
			"LinkToStartContainer: " + DefaultLocalhost + content.container.ContainerID,
	)

	if err := e.email.Send(); err != nil {
		return err
	}
	e.Unlock()

	return nil
}
