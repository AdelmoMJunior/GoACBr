package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
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
	compRepo  repository.CompanyRepository
	certRepo  repository.CertificateRepository
	invRepo   repository.InvoiceRepository
	pool      *acbr.HandlePool
	cryptoSvc crypto.Service
}

func NewEventService(
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	invRepo repository.InvoiceRepository,
	pool *acbr.HandlePool,
	cryptoSvc crypto.Service,
) EventService {
	return &eventService{
		compRepo:  compRepo,
		certRepo:  certRepo,
		invRepo:   invRepo,
		pool:      pool,
		cryptoSvc: cryptoSvc,
	}
}

func (s *eventService) Cancel(ctx context.Context, companyID uuid.UUID, req *dto.CancelRequest) (*dto.EventResponse, error) {
	hd, err := s.pool.GetHandle(ctx, companyID)
	if err != nil {
		return nil, err
	}
	defer s.pool.ReleaseHandle(hd)

	if hd.ConfiguredFor != companyID {
		if err := configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo, s.cryptoSvc); err != nil {
			return nil, err
		}
	}

	// 1. Gerar INI de Evento para Cancelamento
	eventINI := fmt.Sprintf(`[EVENTO]
idLote=%d
[EVENTO001]
chNFe=%s
tpEvento=110111
nSeqEvento=1
dhEvento=%s
xJust=%s
CNPJ=%s`, req.Lote, req.Chave, time.Now().Format("02/01/2006 15:04:05"), req.Justificativa, req.CNPJCPF)

	if err := hd.CarregarEventoINI(eventINI); err != nil {
		return nil, err
	}

	respStr, err := hd.EnviarEvento(req.Lote)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)

	proto := extractFromINI(respStr, "", "nProt")

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
		if err := configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo, s.cryptoSvc); err != nil {
			return nil, err
		}
	}

	// 1. Gerar INI de Evento para CCe
	eventINI := fmt.Sprintf(`[EVENTO]
idLote=%d
[EVENTO001]
chNFe=%s
tpEvento=110110
nSeqEvento=%d
dhEvento=%s
xCorr=%s
CNPJ=%s`, req.Lote, req.Chave, req.NSeqEvento, time.Now().Format("02/01/2006 15:04:05"), req.Correcao, req.CNPJCPF)

	if err := hd.CarregarEventoINI(eventINI); err != nil {
		return nil, err
	}

	respStr, err := hd.EnviarEvento(req.Lote)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)
	proto := extractFromINI(respStr, "", "nProt")

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
		Protocolo:  proto,
		DHEvento:   time.Now(),
	}, nil
}

func (s *eventService) Inutilizacao(ctx context.Context, companyID uuid.UUID, req *dto.InutilizacaoRequest) (*dto.InutilizacaoResponse, error) {
	comp, err := s.compRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	hd, err := s.pool.GetHandle(ctx, companyID)
	if err != nil {
		return nil, err
	}
	defer s.pool.ReleaseHandle(hd)

	if hd.ConfiguredFor != companyID {
		if err := configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo, s.cryptoSvc); err != nil {
			return nil, err
		}
	}

	respStr, err := hd.Inutilizar(comp.CNPJ, req.Justificativa, req.Ano, int(req.Modelo), req.Serie, req.NumInicial, req.NumFinal)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)

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
