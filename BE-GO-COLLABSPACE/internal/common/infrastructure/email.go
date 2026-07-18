package infrastructure

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	SendResetPasswordEmail(toEmail string, resetURL string) error
}

type SMTPEmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPEmailSender(host string, port int, username, password string, from string) EmailSender {
	return &SMTPEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPEmailSender) SendResetPasswordEmail(toEmail string, resetURL string) error {
	subject := "Subject: Password Reset Request\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password. It expires in 15 minutes.</p>
		<a href="%s">Reset Password</a>
	`, resetURL)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, msg)
}
