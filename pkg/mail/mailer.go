package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailSender interface {
	SendMail(to, otpCode string) error
}

type smtpMailer struct {
	Host	string
	Port	string
	User	string
	Pass	string
}

func NewSMTPMailer() EmailSender {
	return &smtpMailer{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
	}
}

func (u *smtpMailer) SendMail(to, otpCode string) error {
	auth := smtp.PlainAuth("", u.User, u.Pass, u.Host)
	msg := fmt.Sprintf("To: %s\r\nSubject: OTP Reset Password\r\n\r\nYour OTP: %s", to, otpCode)

	if err := smtp.SendMail(
		u.Host+":"+u.Port,
		auth,
		u.User,
		[]string{to},
		[]byte(msg),
	); err != nil {
		return err
	}

	return nil
} 