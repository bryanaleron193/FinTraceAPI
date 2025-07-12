package services

import (
	"crypto/tls"
	"log"
	"simple-gin-backend/internal/config"

	"gopkg.in/gomail.v2"
)

// EmailService struct to hold the email configuration
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// NewEmailService initializes a new EmailService
func NewEmailService() *EmailService {

	return &EmailService{
		Host:     config.AppConfig.SMTPHost,
		Port:     config.AppConfig.SMTPPort,
		Username: config.AppConfig.SMTPUser,
		Password: config.AppConfig.SMTPPassword,
		From:     config.AppConfig.EmailFromEmail,
		UseTLS:   config.AppConfig.UseMailTLS,
	}
}

// SendEmail sends an email using the EmailService
func (es *EmailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	// Set the sender and recipient
	m.SetHeader("From", es.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Dial the SMTP server
	d := gomail.NewDialer(es.Host, es.Port, es.Username, es.Password)

	// Disable certificate verification for development purposes
	if es.UseTLS && config.AppConfig.IsDevelopmentEnv {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification
		}
	} else {
		d.TLSConfig = nil
	}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
