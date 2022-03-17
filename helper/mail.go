package helper

import (
	"fmt"
	"net/smtp"
	"strings"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Team Jagokan <jagokan.id@gmail.com>"
const CONFIG_AUTH_EMAIL = "jagokan.id@gmail.com"
const CONFIG_AUTH_PASSWORD = "Jagokan@321!"

func SendMail(to []string, subject string, body string) error {
	cc := []string{}
	err := sendMailGmail(to, cc, subject, body)
	if err != nil {
		logErorr(err)
		return err
	}
	return nil
}

func sendMailGmail(to []string, cc []string, subject, message string) error {
	body := "From: " + CONFIG_SENDER_NAME + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_AUTH_EMAIL, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
