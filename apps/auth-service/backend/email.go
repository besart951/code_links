package main

import (
	"context"
	"fmt"
	"net/smtp"
)

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

type EmailSender interface {
	Send(ctx context.Context, settings SmtpSettings, message EmailMessage, password string) error
}

type NoopEmailSender struct{}

func (NoopEmailSender) Send(_ context.Context, _ SmtpSettings, _ EmailMessage, _ string) error {
	return nil
}

type SmtpEmailSender struct{}

func (SmtpEmailSender) Send(_ context.Context, settings SmtpSettings, message EmailMessage, password string) error {
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
