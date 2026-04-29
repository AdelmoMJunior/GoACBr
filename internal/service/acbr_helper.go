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

// configureHandleForCompany uses a 3-step approach:
//  1. ConfigGravar — ACBr dumps a valid default INI (all fields have valid types)
//  2. ConfigLer   — Load that valid INI back into the handle
//  3. ConfigGravarValor — Override ONLY the company-specific fields via the API
//
// This avoids the "invalid integer" error from hand-crafted INI files.
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

	// 5. Step 1: Let ACBr generate a valid default INI file
	iniDir := "/tmp/acbr_ini"
	os.MkdirAll(iniDir, 0700)
	iniPath := filepath.Join(iniDir, companyID.String()+".ini")

	if err := hd.ConfigGravar(iniPath); err != nil {
		slog.Error("ConfigGravar failed", "path", iniPath, "error", err)
		return fmt.Errorf("failed to generate default ACBr config: %w", err)
	}

	// 6. Step 2: Load the valid default INI back into the handle
	if err := hd.ConfigLer(iniPath); err != nil {
		slog.Error("ConfigLer failed on default INI", "path", iniPath, "error", err)
		return fmt.Errorf("failed to load default ACBr config: %w", err)
	}
	slog.Debug("Default INI loaded successfully", "path", iniPath)

	// 7. Step 3: Override company-specific values via ConfigGravarValor
	//    Use ONLY the keys proven to work from the ACBr API.

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

	ambiente := strconv.Itoa(int(comp.Ambiente))
	if comp.Ambiente == 0 {
		ambiente = "2"
	}

	// Each entry: {section, key, value}
	// These are the fields we override from the defaults.
	// We log-and-skip on error instead of failing hard, since some keys
	// may or may not exist depending on the ACBr version.
	configs := []struct{ section, key, value string }{
		// Principal
		{"Principal", "TipoResposta", "2"},
		{"Principal", "LogNivel", "3"},
		{"Principal", "LogPath", logPath},

		// DFe — certificate
		{"DFe", "SSLCryptLib", "1"},
		{"DFe", "SSLHttpLib", "3"},
		{"DFe", "SSLXmlSignLib", "4"},
		{"DFe", "ArquivoPFX", pfxPath},
		{"DFe", "Senha", pfxPassword},
		{"DFe", "VerificarValidade", "1"},

		// NFe — paths and schemas (using "NFE" section name for API)
		{"NFE", "PathSchemas", schemasPath},
		{"NFE", "PathSalvar", savePath},
		{"NFE", "Ambiente", ambiente},

		// DANFE
		{"DANFE", "PathPDF", pdfPath},
	}

	for _, c := range configs {
		if err := hd.ConfigGravarValor(c.section, c.key, c.value); err != nil {
			// Log but don't fail — some keys may not exist in all ACBr versions
			slog.Warn("ConfigGravarValor skipped (key may not exist)",
				"section", c.section, "key", c.key, "error", err)
		}
	}

	// Save the final config back to disk (so next handle load can reuse it)
	hd.ConfigGravar(iniPath)

	hd.ConfiguredFor = companyID
	slog.Info("ACBr handle configured successfully", "company_id", companyID)
	return nil
}

// extractFromINI helper to extract fields from ACBr INI response.
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
