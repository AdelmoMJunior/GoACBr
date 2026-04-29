package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

// configureHandleForCompany fetches company config and cert, and applies it to the handle.
func configureHandleForCompany(
	ctx context.Context,
	hd *acbr.Handle,
	companyID uuid.UUID,
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
	cryptoSvc crypto.Service,
) error {
	comp, err := compRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	var pfxPath, password string
	cert, err := certRepo.GetByCompanyID(ctx, companyID)
	if err == nil && cert != nil {
		slog.Debug("Found certificate for company, decrypting password...", "company_id", companyID)
		// Decrypt password
		passBytes, err := cryptoSvc.Decrypt(cert.PFXPasswordEnc)
		if err == nil {
			password = string(passBytes)
		} else {
			slog.Error("Failed to decrypt certificate password", "error", err)
		}

		slog.Debug("Decrypting PFX data...", "company_id", companyID)
		// Decrypt PFX and save to temp file for ACBrLib
		pfxEncData, err := cryptoSvc.Decrypt(string(cert.PFXData))
		if err == nil {
			tmpDir := filepath.Join(os.TempDir(), "goacbr", "certs")
			os.MkdirAll(tmpDir, 0700)
			pfxPath = filepath.Join(tmpDir, companyID.String()+".pfx")
			slog.Debug("Writing PFX to temporary file", "path", pfxPath)
			err = os.WriteFile(pfxPath, pfxEncData, 0600)
			if err != nil {
				slog.Error("Failed to write temporary PFX file", "error", err)
			}
		} else {
			slog.Error("Failed to decrypt PFX data", "error", err)
		}
	} else {
		slog.Warn("No certificate found for company", "company_id", companyID)
	}

	slog.Debug("Applying company configuration to ACBr handle", "company_id", companyID, "uf", comp.UF, "ambiente", comp.Ambiente)

	// Ensure directories exist
	_ = os.MkdirAll("/app/logs/acbr", 0755)
	_ = os.MkdirAll("/app/data/nfe", 0755)

	// Build ACBr Config Map based on the correct sections
	cfg := map[string]map[string]string{
		"Principal": {
			"TipoResposta": "2", // INI
			"LogNivel":     "4", // Paranoico
			"LogPath":      "/app/logs/acbr",
		},
		"DFe": {
			"SSLLib":     "1", // libOpenSSL
			"CryptLib":   "1", // libOpenSSL
			"HttpLib":    "3", // libHttpLibCurl
			"XmlSignLib": "4", // libXmlSec
			"ArquivoPFX": pfxPath,
			"Senha":      password,
		},
		"NFe": {
			"ModeloDF":              "55",
			"VersaoDF":              "4.00",
			"SalvarGer":             "1",
			"PathSalvar":            "/app/data/nfe",
			"AtualizarXMLCancelado": "1",
		},
		"WebService": {
			"UF":       comp.UF,
			"Ambiente": fmt.Sprintf("%d", comp.Ambiente),
			"Visualizar": "0",
			"Salvar":     "1",
		},
	}

	return hd.ApplyCompanyConfig(companyID, cfg)
}

// extractFromINI helper to extract fields from ACBr INI response.
// If section is provided, it searches only within that section.
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
	// ACBr INI usually uses dd/mm/yyyy hh:mm:ss or yyyy-mm-ddThh:mm:ss
	// We try to handle common formats
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
