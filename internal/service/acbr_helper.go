package service

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

// ufToIBGECode maps Brazilian state abbreviations to IBGE numeric codes.
var ufToIBGECode = map[string]int{
	"AC": 12, "AL": 27, "AM": 13, "AP": 16, "BA": 29,
	"CE": 23, "DF": 53, "ES": 32, "GO": 52, "MA": 21,
	"MG": 31, "MS": 50, "MT": 51, "PA": 15, "PB": 25,
	"PE": 26, "PI": 22, "PR": 41, "RJ": 33, "RN": 24,
	"RO": 11, "RR": 14, "RS": 43, "SC": 42, "SE": 28,
	"SP": 35, "TO": 17,
}

// UFToCode converts a UF abbreviation (e.g. "SP") to the IBGE numeric code (35).
func UFToCode(uf string) int {
	if code, ok := ufToIBGECode[strings.ToUpper(uf)]; ok {
		return code
	}
	return 91
}

// configureHandleForCompany configures an ACBr handle for a specific company.
//
// Since handles are initialized with a real INI file path (NFE_Inicializar),
// all config sections ([Principal], [DFe], [NFe], [DANFE]) are registered
// with valid defaults. We just override the company-specific values via
// ConfigGravarValor using the exact key names from the ACBrLib docs:
//   - https://acbr.sourceforge.io/ACBrLib/Geral.html       → [Principal]
//   - https://acbr.sourceforge.io/ACBrLib/DFe.html          → [DFe]
//   - https://acbr.sourceforge.io/ACBrLib/NFe2.html          → [NFe]
func configureHandleForCompany(
	ctx context.Context,
	hd *acbr.Handle,
	pool *acbr.HandlePool,
	companyID uuid.UUID,
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	cryptoSvc crypto.Service,
) error {
	// 1. Fetch company data
	comp, err := compRepo.GetByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("failed to fetch company %s: %w", companyID, err)
	}

	// 2. Fetch certificate
	cert, err := certRepo.GetByCompanyID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company %s has no certificate: %w", companyID, err)
	}

	// 3. Decrypt PFX password
	passBytes, err := cryptoSvc.Decrypt(cert.PFXPasswordEnc)
	if err != nil {
		return fmt.Errorf("failed to decrypt certificate password: %w", err)
	}
	pfxPassword := string(passBytes)

	// 4. Decrypt PFX data and write to temp file
	pfxBytes, err := cryptoSvc.Decrypt(string(cert.PFXData))
	if err != nil {
		return fmt.Errorf("failed to decrypt certificate PFX: %w", err)
	}

	certsDir := "/tmp/acbr_certs"
	os.MkdirAll(certsDir, 0700)
	pfxPath := filepath.Join(certsDir, companyID.String()+".pfx")
	if err := os.WriteFile(pfxPath, pfxBytes, 0600); err != nil {
		return fmt.Errorf("failed to write PFX to temp file: %w", err)
	}

	slog.Info("Configuring ACBr handle for company",
		"company_id", companyID,
		"cnpj", comp.CNPJ,
		"uf", comp.UF,
		"ambiente", comp.Ambiente,
	)

	// 5. Prepare paths
	schemasPath := pool.SchemasPath
	if schemasPath == "" {
		schemasPath = "/app/data/Schemas/NFe"
	}
	logPath := pool.LogPath
	if logPath == "" {
		logPath = "/tmp/acbr_logs"
	}
	os.MkdirAll(logPath, 0755)

	savePath := "/tmp/acbr_nfe/" + companyID.String() + "/"
	pdfPath := "/tmp/acbr_pdf/" + companyID.String() + "/"
	os.MkdirAll(savePath, 0755)
	os.MkdirAll(pdfPath, 0755)

	// Ambiente: 0=taProducao, 1=taHomologacao (ACBr uses 0-indexed)
	ambiente := "1" // default homologação
	if comp.Ambiente == 1 {
		ambiente = "0" // produção
	}

	// 6. Apply company config via ConfigGravarValor
	//    Key names match the ACBrLib documentation exactly.
	configs := map[string]map[string]string{
		// [Principal] — Geral.html
		"Principal": {
			"TipoResposta": "2", // JSON
			"LogNivel":     "3", // Completo
			"LogPath":      logPath,
		},
		// [DFe] — DFe.html
		"DFe": {
			"SSLCryptLib":      "1", // cryOpenSSL
			"SSLHttpLib":       "3", // httpOpenSSL
			"SSLXmlSignLib":    "4", // xsLibXml2
			"UF":               comp.UF,
			"ArquivoPFX":       pfxPath,
			"Senha":            pfxPassword,
			"VerificarValidade": "1",
		},
		// [NFe] — NFe2.html (Configurações da Biblioteca)
		"NFe": {
			"Ambiente":     ambiente,
			"ModeloDF":     "0",    // moNFe (55)
			"VersaoDF":     "3",    // ve400
			"SSLType":      "5",    // LT_TLSv1_2
			"Timeout":      "30000",
			"Tentativas":   "5",
			"IntervaloTentativas": "1000",
			"PathSchemas":  schemasPath,
			"PathSalvar":   savePath,
			"SalvarGer":    "1",
			"SalvarEvento": "1",
			"SalvarApenasNFeProcessadas": "1",
			"NormatizarMunicipios":       "1",
			"ExibirErroSchema":           "1",
			"EmissaoPathNFe":             "1",
		},
		// [DANFE]
		"DANFE": {
			"PathPDF": pdfPath,
		},
	}

	hd.ApplyCompanyConfig(configs)

	// 7. Save the final state back to the handle's INI file
	if hd.IniPath != "" {
		hd.ConfigGravar(hd.IniPath)
	}

	hd.ConfiguredFor = companyID
	slog.Info("ACBr handle configured successfully", "company_id", companyID)
	return nil
}

// extractFromINI helper to extract fields from ACBr INI/JSON response.
func extractFromINI(content, section, key string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	prefix := key + "="
	inSection := section == ""

	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		if l == "" {
			continue
		}

		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			currSection := strings.Trim(l, "[]")
			if section != "" {
				inSection = strings.EqualFold(currSection, section)
			}
			continue
		}

		if inSection && strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return "UNKNOWN"
}

func parseINIInt(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}

func parseINIDecimal(s string) decimal.Decimal {
	s = strings.Replace(s, ",", ".", -1)
	d, _ := decimal.NewFromString(s)
	return d
}

func parseINITime(s string) time.Time {
	formats := []string{
		"2006-01-02T15:04:05-07:00",
		"02/01/2006 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t
		}
	}
	return time.Now()
}

func getSections(content string) []string {
	var sections []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			sections = append(sections, strings.Trim(l, "[]"))
		}
	}
	return sections
}
