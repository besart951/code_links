package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
)

type Message struct {
	To      string
	Subject string
	Body    string
}

type Sender interface {
	Send(ctx context.Context, settings domain.SmtpSettings, message Message, password string) error
}

type NoopSender struct{}

func (NoopSender) Send(_ context.Context, _ domain.SmtpSettings, _ Message, _ string) error {
	return nil
}

type SmtpSender struct{}

func (SmtpSender) Send(_ context.Context, settings domain.SmtpSettings, message Message, password string) error {
	if !settings.Active {
		return nil
	}
	if settings.Host == "" || settings.FromEmail == "" || message.To == "" {
		return fmt.Errorf("smtp settings are incomplete")
	}

	addr := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	auth := smtp.PlainAuth("", settings.Username, password, settings.Host)
	content := []byte("To: " + message.To + "\r\n" +
		"From: " + settings.FromEmail + "\r\n" +
		"Subject: " + message.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n\r\n" +
		message.Body)

	return smtp.SendMail(addr, auth, settings.FromEmail, []string{message.To}, content)
}
