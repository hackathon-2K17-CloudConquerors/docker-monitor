package monitor

type EmailManipulator interface {
	SendEmail(from string, to string) error
}
