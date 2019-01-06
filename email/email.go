package email

import (
	"net/smtp"
	"os"
)

func SendBookRequest(title string, telegram string) error {
	from := os.Getenv("EMAIL")
	pass := os.Getenv("GOOGLE_PASS")
	to := os.Getenv("EMAIL")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: book request: " + title + "\n\n" +
		telegram
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
