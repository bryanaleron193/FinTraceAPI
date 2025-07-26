package services

import (
	"simple-gin-backend/internal/config"

	"github.com/resend/resend-go/v2"
)

func SendEmail(to []string, subject string, body string) error {
	client := resend.NewClient(config.AppConfig.ResendAPIKey)

	params := &resend.SendEmailRequest{
		From:    config.AppConfig.ResendFromName + " <" + config.AppConfig.ResendFromEmail + ">",
		To:      to,
		Bcc:     []string{config.AppConfig.AdminEmail},
		Subject: subject,
		Html:    body,
	}

	_, err := client.Emails.Send(params)

	if err != nil {
		return err
	}

	return nil
}
