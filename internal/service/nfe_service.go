package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/storage"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

type NFeService interface {
	Emit(ctx context.Context, companyID uuid.UUID, req *dto.NFeEmitRequest) (*dto.NFeResponse, error)
	QueryStatus(ctx context.Context, companyID uuid.UUID, req *dto.NFeStatusRequest) (*dto.NFeResponse, error)
}

type nfeService struct {
	compRepo   repository.CompanyRepository
	certRepo   repository.CertificateRepository
	invRepo    repository.InvoiceRepository
	pool       *acbr.HandlePool
	storage    storage.Provider
}

func NewNFeService(
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	invRepo repository.InvoiceRepository,
	pool *acbr.HandlePool,
	storage storage.Provider,
) NFeService {
	return &nfeService{
		compRepo: compRepo,
		certRepo: certRepo,
		invRepo:  invRepo,
		pool:     pool,
		storage:  storage,
	}
}

func (s *nfeService) Emit(ctx context.Context, companyID uuid.UUID, req *dto.NFeEmitRequest) (*dto.NFeResponse, error) {
	// 1. Get Handle
	hd, err := s.pool.GetHandle(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ACBr handle: %w", err)
	}
	defer s.pool.ReleaseHandle(hd)

	// 2. Configure Handle if needed
	if hd.ConfiguredFor != companyID {
		err = configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo)
		if err != nil {
			return nil, err
		}
	}

	// 3. Clear List & Load INI
	if err := hd.LimparLista(); err != nil {
		return nil, err
	}

	if err := hd.CarregarINI(req.INIContent); err != nil {
		return nil, apperror.NewBadRequest("invalid INI content: " + err.Error())
	}

	// 4. Sign and Validate
	if err := hd.Assinar(); err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}
	if err := hd.Validar(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 5. Send to SEFAZ
	respStr, err := hd.Enviar(req.Lote, req.PrintPDF, true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send: %w", err)
	}

	// 6. Parse response (in a real app, parse the INI response from ACBr)
	// Here we mock the parsing for scaffold purposes.
	chave := extractFromINI(respStr, "ChaveDFe")
	status := extractFromINI(respStr, "xMotivo")
	cStat := 100 // Mock

	// 7. Store XML in B2 (Assume we extract XML from ACBr, here we just use respStr as mock)
	xmlKey := fmt.Sprintf("%s/%s/%s-nfe.xml", companyID.String(), time.Now().Format("2006/01"), chave)
	_, _ = s.storage.Upload(ctx, xmlKey, bytes.NewReader([]byte("<xml>mock</xml>")), "application/xml")

	// 8. Save Invoice to DB
	inv := &domain.Invoice{
		ID:        uuid.New(),
		CompanyID: companyID,
		Chave:     chave,
		Modelo:    int16(req.Modelo),
		Status:    status,
		XMLB2Key:  xmlKey,
	}
	_ = s.invRepo.Create(ctx, inv)

	return &dto.NFeResponse{
		Chave:    chave,
		Status:   status,
		CStat:    cStat,
		XMLB2Key: xmlKey,
	}, nil
}

func (s *nfeService) QueryStatus(ctx context.Context, companyID uuid.UUID, req *dto.NFeStatusRequest) (*dto.NFeResponse, error) {
	hd, err := s.pool.GetHandle(ctx, companyID)
	if err != nil {
		return nil, err
	}
	defer s.pool.ReleaseHandle(hd)

	if hd.ConfiguredFor != companyID {
		if err := configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo); err != nil {
			return nil, err
		}
	}

	respStr, err := hd.Consultar(req.Chave)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "xMotivo")
	
	// Update DB status asynchronously or synchronously
	_ = s.invRepo.UpdateStatus(ctx, uuid.Nil, status) // Needs actual invoice ID usually, or lookup by chave

	return &dto.NFeResponse{
		Chave:  req.Chave,
		Status: status,
	}, nil
}

// helper to extract fields from ACBr INI response
func extractFromINI(content, key string) string {
	lines := strings.Split(content, "\n")
	prefix := key + "="
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return "UNKNOWN"
}
