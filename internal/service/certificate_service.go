package service

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/pkcs12"

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
	compRepo  repository.CompanyRepository
	cryptoSvc crypto.Service
}

func NewCertificateService(repo repository.CertificateRepository, compRepo repository.CompanyRepository, cryptoSvc crypto.Service) CertificateService {
	return &certificateService{
		repo:      repo,
		compRepo:  compRepo,
		cryptoSvc: cryptoSvc,
	}
}

func (s *certificateService) Upload(ctx context.Context, companyID uuid.UUID, password string, file io.Reader) (*dto.CertificateResponse, error) {
	pfxData, err := io.ReadAll(file)
	if err != nil {
		return nil, apperror.NewBadRequest("failed to read certificate file")
	}

	// Extract metadata from PFX
	// Note: ACBr generally uses PFX (PKCS#12)
	blocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, apperror.NewBadRequest("invalid certificate password or corrupted file")
	}

	var certX509 *x509.Certificate
	for _, b := range blocks {
		if b.Type == "CERTIFICATE" {
			c, err := x509.ParseCertificate(b.Bytes)
			if err == nil {
				certX509 = c
				break
			}
		}
	}

	if certX509 == nil {
		return nil, apperror.NewBadRequest("no valid certificate found in PFX")
	}

	// Validate that the certificate CNPJ matches the company CNPJ
	comp, err := s.compRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, apperror.NewBadRequest("company not found")
	}

	certCN := certX509.Subject.CommonName
	// Brazilian e-CNPJ certificates have the CNPJ in the CN field (e.g., "EMPRESA:03748056000117")
	cnpjFromCert := extractCNPJFromCN(certCN)
	if cnpjFromCert != "" && cnpjFromCert != comp.CNPJ {
		return nil, apperror.NewBadRequest(
			fmt.Sprintf("certificate CNPJ (%s) does not match company CNPJ (%s)", cnpjFromCert, comp.CNPJ),
		)
	}

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
		PFXData:        []byte(encPfx),
		PFXPasswordEnc: encPassword,
		SubjectCN:      certX509.Subject.CommonName,
		SerialNumber:   certX509.SerialNumber.String(),
		ValidFrom:      certX509.NotBefore,
		ValidUntil:     certX509.NotAfter,
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

// cnpjRegex matches a 14-digit CNPJ sequence.
var cnpjRegex = regexp.MustCompile(`\d{14}`)

// extractCNPJFromCN extracts a CNPJ from the certificate's Subject CN.
// Brazilian e-CNPJ certs typically embed the CNPJ in the CN field,
// either as "EMPRESA:12345678000190" or within the full CN string.
func extractCNPJFromCN(cn string) string {
	match := cnpjRegex.FindString(cn)
	return match
}
