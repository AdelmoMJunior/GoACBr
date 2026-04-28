package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
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
	compRepo  repository.CompanyRepository
	certRepo  repository.CertificateRepository
	invRepo   repository.InvoiceRepository
	pool      *acbr.HandlePool
	storage   storage.Provider
	cryptoSvc crypto.Service
}

func NewNFeService(
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	invRepo repository.InvoiceRepository,
	pool *acbr.HandlePool,
	storage storage.Provider,
	cryptoSvc crypto.Service,
) NFeService {
	return &nfeService{
		compRepo:  compRepo,
		certRepo:  certRepo,
		invRepo:   invRepo,
		pool:      pool,
		storage:   storage,
		cryptoSvc: cryptoSvc,
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
		err = configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo, s.cryptoSvc)
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

	// 6. Get signed XML from ACBr
	xmlContent, err := hd.ObterXml(0)
	if err != nil {
		slog.Warn("Failed to get signed XML from ACBr", "error", err)
	}

	// 7. Parse response
	chave := extractFromINI(respStr, "", "ChaveDFe")
	status := extractFromINI(respStr, "", "xMotivo")
	cStatStr := extractFromINI(respStr, "", "cStat")
	cStat := 0
	fmt.Sscanf(cStatStr, "%d", &cStat)

	// 8. Store XML in B2
	xmlKey := fmt.Sprintf("%s/%s/%s-nfe.xml", companyID.String(), time.Now().Format("2006/01"), chave)
	if xmlContent != "" {
		_, _ = s.storage.Upload(ctx, xmlKey, strings.NewReader(xmlContent), "application/xml")
	}

	// 9. Save Invoice to DB
	inv := &domain.Invoice{
		ID:        uuid.New(),
		CompanyID: companyID,
		Chave:     chave,
		Modelo:    int16(req.Modelo),
		Status:    status,
		XMLB2Key:  xmlKey,
		
		// Identificação
		CNF:       extractFromINI(req.INIContent, "Ide", "cNF"),
		NatOp:     extractFromINI(req.INIContent, "Ide", "natOp"),
		Numero:    parseINIInt(extractFromINI(req.INIContent, "Ide", "nNF")),
		Serie:     parseINIInt(extractFromINI(req.INIContent, "Ide", "serie")),
		DHEmissao: parseINITime(extractFromINI(req.INIContent, "Ide", "dhEmi")),
		TpNF:      int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "tpNF"))),
		IDDest:    int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "idDest"))),
		CMunFG:    extractFromINI(req.INIContent, "Ide", "cMunFG"),
		TpImp:     int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "tpImp"))),
		TpEmis:    int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "tpEmis"))),
		TpAmb:     int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "tpAmb"))),
		FinNFe:    int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "finNFe"))),
		IndFinal:  int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "indFinal"))),
		IndPres:   int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "indPres"))),
		ProcEmi:   int16(parseINIInt(extractFromINI(req.INIContent, "Ide", "procEmi"))),
		VerProc:   extractFromINI(req.INIContent, "Ide", "verProc"),
		
		// Destinatário
		DestCNPJCPF: extractFromINI(req.INIContent, "Dest", "CNPJCPF"),
		DestNome:    extractFromINI(req.INIContent, "Dest", "xNome"),
		DestIE:      extractFromINI(req.INIContent, "Dest", "IE"),
		DestEmail:   extractFromINI(req.INIContent, "Dest", "email"),
		
		// Totais
		TotVBC:        parseINIDecimal(extractFromINI(req.INIContent, "Total", "vBC")),
		TotVICMS:      parseINIDecimal(extractFromINI(req.INIContent, "Total", "vICMS")),
		TotVICMSDeson: parseINIDecimal(extractFromINI(req.INIContent, "Total", "vICMSDeson")),
		TotVFCP:       parseINIDecimal(extractFromINI(req.INIContent, "Total", "vFCP")),
		TotVBCST:      parseINIDecimal(extractFromINI(req.INIContent, "Total", "vBCST")),
		TotVST:        parseINIDecimal(extractFromINI(req.INIContent, "Total", "vST")),
		TotVProd:      parseINIDecimal(extractFromINI(req.INIContent, "Total", "vProd")),
		TotVFrete:     parseINIDecimal(extractFromINI(req.INIContent, "Total", "vFrete")),
		TotVSeg:       parseINIDecimal(extractFromINI(req.INIContent, "Total", "vSeg")),
		TotVDesc:      parseINIDecimal(extractFromINI(req.INIContent, "Total", "vDesc")),
		TotVNF:        parseINIDecimal(extractFromINI(req.INIContent, "Total", "vNF")),
		
		// Protocolo e Extras
		Protocolo:  extractFromINI(respStr, "", "nProt"),
		InfAdFisco: extractFromINI(req.INIContent, "InfAdic", "infAdFisco"),
		InfCpl:     extractFromINI(req.INIContent, "InfAdic", "infCpl"),
	}
	
	dhRecStr := extractFromINI(respStr, "", "dhRecbto")
	if dhRecStr != "UNKNOWN" {
		t := parseINITime(dhRecStr)
		inv.DHRecebimento = &t
	}

	// 10. Generate PDF if requested
	if req.PrintPDF {
		if err := hd.ImprimirPDF(); err == nil {
			path, _ := hd.ObterCaminhoGerado()
			if path != "" {
				file, err := os.Open(path)
				if err == nil {
					pdfKey := fmt.Sprintf("%s/%s/%s-nfe.pdf", companyID.String(), time.Now().Format("2006/01"), chave)
					_, err = s.storage.Upload(ctx, pdfKey, file, "application/pdf")
					file.Close()
					if err == nil {
						inv.PDFB2Key = pdfKey
						_ = os.Remove(path) // Clean up temporary PDF
					}
				}
			}
		}
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
		if err := configureHandleForCompany(ctx, hd, companyID, s.compRepo, s.certRepo, s.cryptoSvc); err != nil {
			return nil, err
		}
	}

	respStr, err := hd.Consultar(req.Chave)
	if err != nil {
		return nil, err
	}

	status := extractFromINI(respStr, "", "xMotivo")
	
	// Update DB status asynchronously or synchronously
	_ = s.invRepo.UpdateStatus(ctx, uuid.Nil, status) // Needs actual invoice ID usually, or lookup by chave

	return &dto.NFeResponse{
		Chave:  req.Chave,
		Status: status,
	}, nil
}
