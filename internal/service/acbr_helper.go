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
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
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

// configureHandleForCompany fetches real company/cert data from DB,
// generates a complete ACBrLib INI file, and loads it into the handle.
// This approach is more reliable than ConfigGravarValor because the INI
// format section/key names are well-documented and stable.
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

	// 5. Generate the complete INI content and write to file
	iniContent := generateACBrINI(comp, pfxPath, pfxPassword, pool)

	iniDir := "/tmp/acbr_ini"
	if err := os.MkdirAll(iniDir, 0700); err != nil {
		return fmt.Errorf("failed to create ini dir: %w", err)
	}
	iniPath := filepath.Join(iniDir, companyID.String()+".ini")
	if err := os.WriteFile(iniPath, []byte(iniContent), 0600); err != nil {
		return fmt.Errorf("failed to write INI file: %w", err)
	}

	// 6. Load the INI into the handle
	if err := hd.ConfigLer(iniPath); err != nil {
		return fmt.Errorf("failed to load ACBr config: %w", err)
	}

	// Ensure output directories exist
	os.MkdirAll("/tmp/acbr_nfe/"+companyID.String(), 0755)
	os.MkdirAll("/tmp/acbr_pdf/"+companyID.String(), 0755)

	hd.ConfiguredFor = companyID
	slog.Info("ACBr handle configured successfully", "company_id", companyID)
	return nil
}

// generateACBrINI creates a complete ACBrLib.ini content string for a company.
// Section and key names here match the INI file format exactly.
func generateACBrINI(comp *domain.Company, pfxPath, pfxPassword string, pool *acbr.HandlePool) string {
	schemasPath := pool.SchemasPath
	if schemasPath == "" {
		schemasPath = "/app/data/Schemas/NFe"
	}

	logPath := pool.LogPath
	if logPath == "" {
		logPath = "/tmp/acbr_logs"
	}
	os.MkdirAll(logPath, 0755)

	savePath := "/tmp/acbr_nfe/" + comp.ID.String() + "/"
	pdfPath := "/tmp/acbr_pdf/" + comp.ID.String() + "/"

	ambiente := strconv.Itoa(int(comp.Ambiente))
	if comp.Ambiente == 0 {
		ambiente = "2" // Default to homologação
	}

	var b strings.Builder

	// [Principal]
	b.WriteString("[Principal]\n")
	b.WriteString("TipoResposta=2\n")      // JSON
	b.WriteString("CodificacaoResposta=0\n") // UTF-8
	b.WriteString("LogNivel=3\n")
	b.WriteString("LogPath=" + logPath + "\n")
	b.WriteString("\n")

	// [DFe] — SSL and Certificate
	b.WriteString("[DFe]\n")
	b.WriteString("SSLLib=1\n")       // OpenSSL
	b.WriteString("CryptLib=1\n")     // OpenSSL
	b.WriteString("HttpLib=3\n")      // libcurl/OpenSSL
	b.WriteString("XmlSignLib=4\n")   // libxml2
	b.WriteString("SSLType=5\n")      // TLS 1.2
	b.WriteString("ArquivoPFX=" + pfxPath + "\n")
	b.WriteString("Senha=" + pfxPassword + "\n")
	b.WriteString("VerificarValidade=1\n")
	b.WriteString("AguardarConsultaRet=0\n")
	b.WriteString("IntervaloTentativas=1000\n")
	b.WriteString("Tentativas=5\n")
	b.WriteString("Timeout=15000\n")
	b.WriteString("QuebradeLinha=|\n")
	b.WriteString("\n")

	// [NFe] — Model, version, behavior
	b.WriteString("[NFe]\n")
	b.WriteString("FormaEmissao=0\n")
	b.WriteString("ModeloDF=55\n")
	b.WriteString("VersaoDF=4.00\n")
	b.WriteString("SalvarTXT=0\n")
	b.WriteString("SalvarXML=1\n")
	b.WriteString("SalvarEvento=1\n")
	b.WriteString("SalvarApenasNFeProcessadas=1\n")
	b.WriteString("EmissaoPathNFe=1\n")
	b.WriteString("NormatizarMunicipios=1\n")
	b.WriteString("ExibirErroSchema=1\n")
	b.WriteString("SalvarGer=1\n")
	b.WriteString("AtualizarXMLCancelado=1\n")
	if comp.CSCID != "" {
		b.WriteString("IdCSC=" + comp.CSCID + "\n")
	}
	if comp.CSCToken != "" {
		b.WriteString("CSC=" + comp.CSCToken + "\n")
	}
	b.WriteString("\n")

	// [WebService] — UF and environment
	b.WriteString("[WebService]\n")
	b.WriteString("UF=" + comp.UF + "\n")
	b.WriteString("Ambiente=" + ambiente + "\n")
	b.WriteString("Visualizar=0\n")
	b.WriteString("Salvar=1\n")
	b.WriteString("AjustaAguarda=1\n")
	b.WriteString("Aguardar=5000\n")
	b.WriteString("Tentativas=5\n")
	b.WriteString("IntervaloTentativas=1000\n")
	b.WriteString("\n")

	// [Arquivos] — Storage paths and schemas
	b.WriteString("[Arquivos]\n")
	b.WriteString("Salvar=1\n")
	b.WriteString("SepararPorMes=1\n")
	b.WriteString("SepararPorCNPJ=1\n")
	b.WriteString("SepararPorModelo=1\n")
	b.WriteString("AdicionarLiteral=1\n")
	b.WriteString("EmissaoPathNFe=1\n")
	b.WriteString("PathSalvar=" + savePath + "\n")
	b.WriteString("PathSchemas=" + schemasPath + "\n")
	b.WriteString("\n")

	// [DANFE] — PDF output
	b.WriteString("[DANFE]\n")
	b.WriteString("PathPDF=" + pdfPath + "\n")
	b.WriteString("\n")

	return b.String()
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
