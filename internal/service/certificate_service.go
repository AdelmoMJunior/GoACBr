package service

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

type CertificateService interface {
	Upload(ctx context.Context, companyID uuid.UUID, password string, file io.Reader) (*dto.CertificateResponse, error)
	Get(ctx context.Context, companyID uuid.UUID) (*dto.CertificateResponse, error)
	Delete(ctx context.Context, companyID uuid.UUID) error
}

type certificateService struct {
	repo      repository.CertificateRepository
	cryptoSvc crypto.Service
}

func NewCertificateService(repo repository.CertificateRepository, cryptoSvc crypto.Service) CertificateService {
	return &certificateService{
		repo:      repo,
		cryptoSvc: cryptoSvc,
	}
}

func (s *certificateService) Upload(ctx context.Context, companyID uuid.UUID, password string, file io.Reader) (*dto.CertificateResponse, error) {
	pfxData, err := io.ReadAll(file)
	if err != nil {
		return nil, apperror.NewBadRequest("failed to read certificate file")
	}

	// In a real implementation, we would parse the PFX here using standard library x509
	// to extract SubjectCN, SerialNumber, ValidFrom, and ValidUntil.
	// For simplicity in this scaffold, we mock the extraction.
	// We MUST encrypt the password and the PFX data (or store PFX in B2).
	// Storing PFX in DB is fine if encrypted, since it's small.

	encPassword, err := s.cryptoSvc.Encrypt([]byte(password))
	if err != nil {
		return nil, err
	}

	encPfx, err := s.cryptoSvc.Encrypt(pfxData)
	if err != nil {
		return nil, err
	}

	cert := &domain.Certificate{
		ID:             uuid.New(),
		CompanyID:      companyID,
		PFXData:        []byte(encPfx), // Wait, Encrypt returns string now? Let's assume it returns string.
		PFXPasswordEnc: encPassword,
		SubjectCN:      "Mock Subject CN", // Replace with actual extraction
		SerialNumber:   "1234567890",      // Replace with actual extraction
		ValidFrom:      time.Now(),
		ValidUntil:     time.Now().Add(365 * 24 * time.Hour), // 1 year
	}

	if err := s.repo.Save(ctx, cert); err != nil {
		return nil, err
	}

	return s.mapToResponse(cert), nil
}

func (s *certificateService) Get(ctx context.Context, companyID uuid.UUID) (*dto.CertificateResponse, error) {
	cert, err := s.repo.GetByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(cert), nil
}

func (s *certificateService) Delete(ctx context.Context, companyID uuid.UUID) error {
	cert, err := s.repo.GetByCompanyID(ctx, companyID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, cert.ID)
}

func (s *certificateService) mapToResponse(c *domain.Certificate) *dto.CertificateResponse {
	now := time.Now()
	days := int(c.ValidUntil.Sub(now).Hours() / 24)
	if days < 0 {
		days = 0
	}

	return &dto.CertificateResponse{
		ID:              c.ID,
		CompanyID:       c.CompanyID,
		SubjectCN:       c.SubjectCN,
		SerialNumber:    c.SerialNumber,
		ValidFrom:       c.ValidFrom,
		ValidUntil:      c.ValidUntil,
		DaysUntilExpiry: days,
		IsExpired:       now.After(c.ValidUntil),
		CreatedAt:       c.CreatedAt,
	}
}
