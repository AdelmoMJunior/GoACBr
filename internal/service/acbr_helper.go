package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

// configureHandleForCompany fetches company config and cert, and applies it to the handle.
func configureHandleForCompany(
	ctx context.Context,
	hd *acbr.Handle,
	companyID uuid.UUID,
	compRepo repository.CompanyRepository,
	certRepo repository.CertificateRepository,
) error {
	comp, err := compRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	// Fetch certificate details (mock password extraction for now)
	_, err = certRepo.GetByCompanyID(ctx, companyID)
	if err != nil {
		// Log warning, might still work if no cert needed for specific call
	}

	// Build ACBr Config Map
	cfg := map[string]map[string]string{
		"Principal": {
			"TipoResposta": "2", // INI
			"LogNivel":     "4", // Paranoico
			"LogPath":      "/tmp/acbr",
		},
		"DFe": {
			"ArquivoPFX": "mock_path_to_pfx_or_buffer",
			"Senha":      "mock_senha", // Need to decrypt from cert repo
			"UF":         comp.UF,
		},
		"NFe": {
			"Ambiente":           fmt.Sprintf("%d", comp.Ambiente),
			"SalvarGer":          "1",
			"PathSalvar":         "/tmp/nfe",
			"AtualizarXMLCancelado": "1",
		},
	}

	return hd.ApplyCompanyConfig(companyID, cfg)
}
