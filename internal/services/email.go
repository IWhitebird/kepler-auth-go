package services

import (
	"fmt"
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/models"
	"net/smtp"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendEmail(req *models.EmailRequest) error {
	if s.cfg.Email.SMTPUser == "" || s.cfg.Email.SMTPPassword == "" {
		fmt.Printf("Email configuration not complete, would send: %+v\n", req)
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUser, s.cfg.Email.SMTPPassword, s.cfg.Email.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.Email.FromEmail, req.To, req.Subject, req.Body)

	err := smtp.SendMail(
		s.cfg.Email.SMTPHost+":"+s.cfg.Email.SMTPPort,
		auth,
		s.cfg.Email.FromEmail,
		[]string{req.To},
		[]byte(msg),
	)

	return err
}
