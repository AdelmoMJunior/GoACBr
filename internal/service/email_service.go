package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

type EmailService interface {
	SendDanfe(ctx context.Context, companyID uuid.UUID, to []string, subject, body string, pdfBytes []byte) error
}

type emailService struct {
	compRepo  repository.CompanyRepository
	cryptoSvc crypto.Service
}

func NewEmailService(compRepo repository.CompanyRepository, cryptoSvc crypto.Service) EmailService {
	return &emailService{
		compRepo:  compRepo,
		cryptoSvc: cryptoSvc,
	}
}

func (s *emailService) SendDanfe(ctx context.Context, companyID uuid.UUID, to []string, subject, body string, pdfBytes []byte) error {
	comp, err := s.compRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	if comp.SMTPHost == "" {
		return fmt.Errorf("company %s does not have SMTP configured", companyID)
	}

	passBytes, err := s.cryptoSvc.Decrypt(comp.SMTPPasswordEnc)
	if err != nil {
		return fmt.Errorf("failed to decrypt smtp password: %w", err)
	}

	auth := smtp.PlainAuth("", comp.SMTPUser, string(passBytes), comp.SMTPHost)
	
	addr := fmt.Sprintf("%s:%d", comp.SMTPHost, comp.SMTPPort)

	// Build the email with attachment (simplified for scaffold)
	// In a real app we use a library like gopkg.in/gomail.v2 or similar
	
	slog.Info("Sending email via SMTP", "host", comp.SMTPHost, "to", to)

	// Very simple mock sending
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to[0], subject, body))
	err = smtp.SendMail(addr, auth, comp.SMTPFrom, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
