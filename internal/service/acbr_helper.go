package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

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
		// Decrypt password
		passBytes, err := cryptoSvc.Decrypt(cert.PFXPasswordEnc)
		if err == nil {
			password = string(passBytes)
		}

		// Decrypt PFX and save to temp file for ACBrLib
		pfxEncData, err := cryptoSvc.Decrypt(string(cert.PFXData))
		if err == nil {
			tmpDir := filepath.Join(os.TempDir(), "goacbr", "certs")
			os.MkdirAll(tmpDir, 0700)
			pfxPath = filepath.Join(tmpDir, companyID.String()+".pfx")
			_ = os.WriteFile(pfxPath, pfxEncData, 0600)
		}
	}

	// Build ACBr Config Map
	cfg := map[string]map[string]string{
		"Principal": {
			"TipoResposta": "2", // INI
			"LogNivel":     "4", // Paranoico
			"LogPath":      "/app/logs/acbr",
		},
		"DFe": {
			"ArquivoPFX": pfxPath,
			"Senha":      password,
			"UF":         comp.UF,
		},
		"NFe": {
			"Ambiente":              fmt.Sprintf("%d", comp.Ambiente),
			"SalvarGer":             "1",
			"PathSalvar":            "/app/data/nfe",
			"AtualizarXMLCancelado": "1",
		},
	}

	return hd.ApplyCompanyConfig(companyID, cfg)
}

// extractFromINI helper to extract fields from ACBr INI response.
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
