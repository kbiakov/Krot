package push_service

import (
	"net/smtp"
	"os"
)

const SmtpServerUri = "smtp.gmail.com:587"

func SendEmail(subsID string, to string, text string) error {
	emailAuth := getEmailAuthInstance()

	msg := "From: " + emailAuth.From + "\n" +
		"To: " + to + "\n" +
		"Subject: Krot #" + subsID + "\n\n" +
		text

	err := smtp.SendMail(
		SmtpServerUri,
		emailAuth.Auth,
		emailAuth.From,
		[]string{to},
		[]byte(msg),
	)

	return err
}

// - Singleton instance

type EmailAuth struct {
	Auth     smtp.Auth
	From     string
	Password string
}

var emailAuthInstance *EmailAuth

func getEmailAuthInstance() *EmailAuth {
	once.Do(func() {
		emailAuthInstance = newEmailAuthService()
	})

	return emailAuthInstance
}

func newEmailAuthService() *EmailAuth {
	from := os.Getenv("SUPERUSER_EMAIL")
	password := os.Getenv("SUPERUSER_PASSWORD")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return &EmailAuth{
		From:     from,
		Password: password,
		Auth:     auth,
	}
}
