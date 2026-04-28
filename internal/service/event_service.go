package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

type EventService interface {
	Cancel(ctx context.Context, companyID uuid.UUID, req *dto.CancelRequest) (*dto.EventResponse, error)
	CCe(ctx context.Context, companyID uuid.UUID, req *dto.CCeRequest) (*dto.EventResponse, error)
	Inutilizacao(ctx context.Context, companyID uuid.UUID, req *dto.InutilizacaoRequest) (*dto.InutilizacaoResponse, error)
}

type eventService struct {
	compRepo repository.CompanyRepository
	certRepo repository.CertificateRepository
	invRepo  repository.InvoiceRepository
	pool     *acbr.HandlePool
}

func NewEventService(
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	invRepo repository.InvoiceRepository,
	pool *acbr.HandlePool,
) EventService {
	return &eventService{
		compRepo: compRepo,
		certRepo: certRepo,
		invRepo:  invRepo,
		pool:     pool,
	}
}

func (s *eventService) Cancel(ctx context.Context, companyID uuid.UUID, req *dto.CancelRequest) (*dto.EventResponse, error) {
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

	respStr, err := hd.Cancelar(req.Chave, req.Justificativa, req.CNPJCPF, req.Lote)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "xMotivo")
	cStat := 101 // Cancelado (mock parsing)

	// In a real app, parse protocol and other details from response
	proto := extractFromINI(respStr, "nProt")

	// Store Event in DB
	event := &domain.InvoiceEvent{
		ID:            uuid.New(),
		CompanyID:     companyID,
		ChaveNFe:      req.Chave,
		TpEvento:      "110111", // Cancelamento
		NSeqEvento:    1,
		DHEvento:      time.Now(),
		Protocolo:     proto,
		XJust:         req.Justificativa,
		DestCNPJCPF:   req.CNPJCPF,
	}
	
	// Check if invoice exists and link
	inv, _ := s.invRepo.GetByChave(ctx, req.Chave)
	if inv != nil {
		event.InvoiceID = &inv.ID
		_ = s.invRepo.UpdateStatus(ctx, inv.ID, status)
	}

	_ = s.invRepo.CreateEvent(ctx, event)

	return &dto.EventResponse{
		Chave:      req.Chave,
		TpEvento:   "110111",
		NSeqEvento: 1,
		CStat:      cStat,
		XMotivo:    status,
		Protocolo:  proto,
		DHEvento:   time.Now(),
	}, nil
}

func (s *eventService) CCe(ctx context.Context, companyID uuid.UUID, req *dto.CCeRequest) (*dto.EventResponse, error) {
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

	// For CCe, ACBr uses NFE_CartaCorrecao
	// We didn't bind it in nfe.go yet, but let's mock it using ACBr methods if it was bound
	// respStr := "mock_cce_response" // err := hd.CartaCorrecao(req.Chave, req.Correcao, req.CNPJCPF, req.Lote)

	status := "Carta de Correcao Registrada"
	cStat := 135

	event := &domain.InvoiceEvent{
		ID:            uuid.New(),
		CompanyID:     companyID,
		ChaveNFe:      req.Chave,
		TpEvento:      "110110", // CCe
		NSeqEvento:    req.NSeqEvento,
		DHEvento:      time.Now(),
		XCorrecao:     req.Correcao,
		DestCNPJCPF:   req.CNPJCPF,
	}
	
	inv, _ := s.invRepo.GetByChave(ctx, req.Chave)
	if inv != nil {
		event.InvoiceID = &inv.ID
	}

	_ = s.invRepo.CreateEvent(ctx, event)

	return &dto.EventResponse{
		Chave:      req.Chave,
		TpEvento:   "110110",
		NSeqEvento: req.NSeqEvento,
		CStat:      cStat,
		XMotivo:    status,
		DHEvento:   time.Now(),
	}, nil
}

func (s *eventService) Inutilizacao(ctx context.Context, companyID uuid.UUID, req *dto.InutilizacaoRequest) (*dto.InutilizacaoResponse, error) {
	// similar logic with ACBr NFE_Inutilizar
	cStat := 102
	status := "Inutilizacao de numero homologado"

	inut := &domain.InvoiceInutilizacao{
		ID:            uuid.New(),
		CompanyID:     companyID,
		Ano:           req.Ano,
		Modelo:        int16(req.Modelo),
		Serie:         req.Serie,
		NumInicial:    req.NumInicial,
		NumFinal:      req.NumFinal,
		Justificativa: req.Justificativa,
		Status:        status,
	}
	_ = s.invRepo.CreateInutilizacao(ctx, inut)

	return &dto.InutilizacaoResponse{
		CStat:   cStat,
		XMotivo: status,
	}, nil
}
