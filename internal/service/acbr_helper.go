package service

import (
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
// Returns 91 (SEFAZ Nacional) if the UF is unknown.
func UFToCode(uf string) int {
	if code, ok := ufToIBGECode[strings.ToUpper(uf)]; ok {
		return code
	}
	return 91
}

// configureHandleForCompany fetches real company/cert data from DB and
// applies it to the ACBr handle. This is the single source of truth for
// all ACBr INI configuration.
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
	if err := os.MkdirAll(certsDir, 0700); err != nil {
		return fmt.Errorf("failed to create certs dir: %w", err)
	}
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

	// 5. Apply all configuration sections
	// [Principal]
	hd.ConfigGravarValor("Principal", "TipoResposta", "2") // 2 = JSON
	hd.ConfigGravarValor("Principal", "LogNivel", "3")
	if pool.LogPath != "" {
		hd.ConfigGravarValor("Principal", "LogPath", pool.LogPath)
	}

	// [DFe] — SSL / Certificate
	hd.ConfigGravarValor("DFe", "SSLCryptLib", "1")    // OpenSSL
	hd.ConfigGravarValor("DFe", "SSLHttpLib", "3")     // WinHTTP/OpenSSL
	hd.ConfigGravarValor("DFe", "SSLXmlSignLib", "4")  // libxml2
	hd.ConfigGravarValor("DFe", "SSLType", "5")        // TLS 1.2
	hd.ConfigGravarValor("DFe", "ArquivoPFX", pfxPath)
	hd.ConfigGravarValor("DFe", "Senha", pfxPassword)
	hd.ConfigGravarValor("DFe", "VerificarValidade", "1")
	hd.ConfigGravarValor("DFe", "Timeout", "15000")
	hd.ConfigGravarValor("DFe", "Tentativas", "5")
	hd.ConfigGravarValor("DFe", "IntervaloTentativas", "1000")

	// [WebService] — UF and environment from company data
	hd.ConfigGravarValor("WebService", "UF", comp.UF)
	hd.ConfigGravarValor("WebService", "Ambiente", strconv.Itoa(int(comp.Ambiente)))
	hd.ConfigGravarValor("WebService", "Salvar", "1")
	hd.ConfigGravarValor("WebService", "AjustaAguarda", "1")
	hd.ConfigGravarValor("WebService", "Aguardar", "5000")
	hd.ConfigGravarValor("WebService", "Tentativas", "5")
	hd.ConfigGravarValor("WebService", "IntervaloTentativas", "1000")

	// [NFe] — Model, version, schema path
	hd.ConfigGravarValor("NFe", "ModeloDF", "55")
	hd.ConfigGravarValor("NFe", "VersaoDF", "4.00")
	hd.ConfigGravarValor("NFe", "SalvarXML", "1")
	hd.ConfigGravarValor("NFe", "SalvarEvento", "1")
	hd.ConfigGravarValor("NFe", "SalvarApenasNFeProcessadas", "1")
	hd.ConfigGravarValor("NFe", "NormatizarMunicipios", "1")
	hd.ConfigGravarValor("NFe", "ExibirErroSchema", "1")

	// CSC/IdCSC for NFCe
	if comp.CSCID != "" {
		hd.ConfigGravarValor("NFe", "IdCSC", comp.CSCID)
	}
	if comp.CSCToken != "" {
		hd.ConfigGravarValor("NFe", "CSC", comp.CSCToken)
	}

	// [Arquivos] — Save paths and schemas
	hd.ConfigGravarValor("Arquivos", "Salvar", "1")
	hd.ConfigGravarValor("Arquivos", "SepararPorMes", "1")
	hd.ConfigGravarValor("Arquivos", "SepararPorCNPJ", "1")
	hd.ConfigGravarValor("Arquivos", "SepararPorModelo", "1")
	hd.ConfigGravarValor("Arquivos", "AdicionarLiteral", "1")
	hd.ConfigGravarValor("Arquivos", "EmissaoPathNFe", "1")
	hd.ConfigGravarValor("Arquivos", "PathSalvar", "/tmp/acbr_nfe/"+companyID.String()+"/")
	if pool.SchemasPath != "" {
		hd.ConfigGravarValor("Arquivos", "PathSchemas", pool.SchemasPath)
	}

	// [DANFE] — PDF output path
	hd.ConfigGravarValor("DANFE", "PathPDF", "/tmp/acbr_pdf/"+companyID.String()+"/")

	// Ensure output directories exist
	os.MkdirAll("/tmp/acbr_nfe/"+companyID.String(), 0755)
	os.MkdirAll("/tmp/acbr_pdf/"+companyID.String(), 0755)

	hd.ConfiguredFor = companyID
	slog.Info("ACBr handle configured successfully", "company_id", companyID)
	return nil
}

// extractFromINI helper to extract fields from ACBr INI response.
func extractFromINI(content, section, key string) string {
	lines := strings.Split(content, "\n")
	prefix := key + "="
	inSection := section == ""

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		// Check for section
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
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			sections = append(sections, strings.Trim(l, "[]"))
		}
	}
	return sections
}
