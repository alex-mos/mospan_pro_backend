package email

import (
	"github.com/alex-mos/mospan_pro_backend/config"
	"net/smtp"
)

func SendBookRequest(title string, telegram string) error {
	from := config.GetConfig().Email
	pass := config.GetConfig().GooglePassword
	to := config.GetConfig().Email

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
