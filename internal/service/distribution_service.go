package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

type DistributionService interface {
	QueryByUltNSU(ctx context.Context, companyID uuid.UUID, ultNSU string) (*dto.DistributionQueryResponse, error)
	QueryByNSU(ctx context.Context, companyID uuid.UUID, nsu string) (*dto.DistributionQueryResponse, error)
	GetControl(ctx context.Context, companyID uuid.UUID) (*dto.DistributionControlResponse, error)
}

type distributionService struct {
	compRepo  repository.CompanyRepository
	certRepo  repository.CertificateRepository
	distRepo  repository.DistributionRepository
	pool      *acbr.HandlePool
	cryptoSvc crypto.Service
}

func NewDistributionService(
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	distRepo repository.DistributionRepository,
	pool *acbr.HandlePool,
	cryptoSvc crypto.Service,
) DistributionService {
	return &distributionService{
		compRepo:  compRepo,
		certRepo:  certRepo,
		distRepo:  distRepo,
		pool:      pool,
		cryptoSvc: cryptoSvc,
	}
}

func (s *distributionService) QueryByUltNSU(ctx context.Context, companyID uuid.UUID, ultNSU string) (*dto.DistributionQueryResponse, error) {
	comp, err := s.compRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	// Verify cooldown
	ctrl, err := s.distRepo.GetControl(ctx, companyID)
	if err == nil && ctrl.LastQueryAt != nil {
		if time.Since(*ctrl.LastQueryAt) < 1*time.Hour && ctrl.MaxNSU == ctrl.LastNSU {
			// SEFAZ rejects queries within 1 hr if there are no new NSUs
			return nil, apperror.NewTooManyRequests("SEFAZ cooldown active. Try again later.")
		}
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

	// AcUFAutor is usually the UF code of the company or SEFAZ Nacional (91)
	ufAutor := 91

	respStr, err := hd.DistribuicaoDFePorUltNSU(ufAutor, comp.CNPJ, ultNSU)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)

	maxNsu := extractFromINI(respStr, "", "maxNSU")
	ultNsuRet := extractFromINI(respStr, "", "ultNSU")

	if maxNsu == "" || maxNsu == "UNKNOWN" {
		maxNsu = ultNSU
		ultNsuRet = ultNSU
	}

	now := time.Now()
	newCtrl := &domain.DistributionControl{
		CompanyID:   companyID,
		LastNSU:     ultNsuRet,
		MaxNSU:      maxNsu,
		LastQueryAt: &now,
		Status:      "idle",
	}
	_ = s.distRepo.UpsertControl(ctx, newCtrl)

	// 10. Extract docs from response and save them
	var docs []dto.DistributionDoc
	sections := getSections(respStr)
	for _, sec := range sections {
		// DFe distribution sections start with ResNFe, ResEvento, ProcNFe, etc.
		if strings.HasPrefix(sec, "ResNFe") || strings.HasPrefix(sec, "ResEvento") ||
			strings.HasPrefix(sec, "ProcNFe") || strings.HasPrefix(sec, "ProcEvento") {

			nsu := extractFromINI(respStr, sec, "NSU")
			chave := extractFromINI(respStr, sec, "chNFe")
			if chave == "UNKNOWN" {
				chave = extractFromINI(respStr, sec, "chNFe")
			}

			doc := dto.DistributionDoc{
				NSU:    nsu,
				Chave:  chave,
				Schema: sec,
			}
			docs = append(docs, doc)

			// Save to database
			dbDoc := &domain.DistributionDocument{
				ID:         uuid.New(),
				CompanyID:  companyID,
				NSU:        nsu,
				ChaveNFe:   chave,
				SchemaType: sec,

				EmitCNPJCPF: extractFromINI(respStr, sec, "CNPJCPF"),
				EmitNome:    extractFromINI(respStr, sec, "xNome"),
				EmitIE:      extractFromINI(respStr, sec, "IE"),
				DestCNPJCPF: extractFromINI(respStr, sec, "dest_CNPJCPF"),

				TpEvento:   extractFromINI(respStr, sec, "tpEvento"),
				DescEvento: extractFromINI(respStr, sec, "xEvento"),
				Protocolo:  extractFromINI(respStr, sec, "nProt"),

				CreatedAt: time.Now(),
			}

			// Conversão de valores com tratamento de ponteiros
			vNF := parseINIDecimal(extractFromINI(respStr, sec, "vNF"))
			if !vNF.IsZero() {
				dbDoc.TotVNF = &vNF
			}

			vICMS := parseINIDecimal(extractFromINI(respStr, sec, "vICMS"))
			if !vICMS.IsZero() {
				dbDoc.TotVICMS = &vICMS
			}

			vST := parseINIDecimal(extractFromINI(respStr, sec, "vST"))
			if !vST.IsZero() {
				dbDoc.TotVST = &vST
			}

			// Datas
			if dhEmiStr := extractFromINI(respStr, sec, "dhEmi"); dhEmiStr != "UNKNOWN" {
				t := parseINITime(dhEmiStr)
				dbDoc.DHEmissao = &t
			}
			if dhRecStr := extractFromINI(respStr, sec, "dhRecbto"); dhRecStr != "UNKNOWN" {
				t := parseINITime(dhRecStr)
				dbDoc.DHRecebimento = &t
			}

			// Campos inteiros
			if modStr := extractFromINI(respStr, sec, "modelo"); modStr != "UNKNOWN" {
				val := int16(parseINIInt(modStr))
				dbDoc.Modelo = &val
			}
			if serStr := extractFromINI(respStr, sec, "serie"); serStr != "UNKNOWN" {
				val := parseINIInt(serStr)
				dbDoc.Serie = &val
			}
			if numStr := extractFromINI(respStr, sec, "numero"); numStr != "UNKNOWN" {
				val := parseINIInt(numStr)
				dbDoc.Numero = &val
			}

			_ = s.distRepo.SaveDocument(ctx, dbDoc)
		}
	}

	return &dto.DistributionQueryResponse{
		CStat:      cStat,
		XMotivo:    status,
		UltNSU:     ultNsuRet,
		MaxNSU:     maxNsu,
		Documentos: docs,
	}, nil
}

func (s *distributionService) QueryByNSU(ctx context.Context, companyID uuid.UUID, nsu string) (*dto.DistributionQueryResponse, error) {
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

	ufAutor := 91
	respStr, err := hd.DistribuicaoDFePorNSU(ufAutor, comp.CNPJ, nsu)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)

	return &dto.DistributionQueryResponse{
		CStat:   cStat,
		XMotivo: status,
	}, nil
}

func (s *distributionService) GetControl(ctx context.Context, companyID uuid.UUID) (*dto.DistributionControlResponse, error) {
	ctrl, err := s.distRepo.GetControl(ctx, companyID)
	if err != nil {
		return nil, err
	}

	var next time.Time
	if ctrl.LastQueryAt != nil {
		next = ctrl.LastQueryAt.Add(1 * time.Hour)
	}

	return &dto.DistributionControlResponse{
		CompanyID:        ctrl.CompanyID.String(),
		LastNSU:          ctrl.LastNSU,
		MaxNSU:           ctrl.MaxNSU,
		LastQueryAt:      ctrl.LastQueryAt,
		IsRunning:        ctrl.IsRunning,
		Status:           ctrl.Status,
		ErrorMessage:     ctrl.ErrorMessage,
		NextAllowedQuery: &next,
	}, nil
}
