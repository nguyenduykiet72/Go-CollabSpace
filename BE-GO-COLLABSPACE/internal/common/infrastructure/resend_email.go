package infrastructure

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendEmailSender struct {
	client    *resend.Client
	fromEmail string
}

func NewResendEmailSender(apiKey string, fromEmail string) *ResendEmailSender {
	client := resend.NewClient(apiKey)
	return &ResendEmailSender{
		client:    client,
		fromEmail: fromEmail,
	}
}

func (s *ResendEmailSender) SendResetPasswordEmail(toEmail string, resetURL string) error {
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: "CollabSpace - Reset your password",
		Html: fmt.Sprintf(`
			<h2>Password Reset Request</h2>
			<p>We received a request to reset your password. Click the link below to set a new one:</p>
			<a href="%s" style="padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a>
			<p>This link will expire in 15 minutes.</p>
			<p>If you didn't request this, you can safely ignore this email.</p>
		`, resetURL),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend failed to send email: %w", err)
	}

	return nil
}
