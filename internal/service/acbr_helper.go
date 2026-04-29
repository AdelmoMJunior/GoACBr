package service

import (
	"context"
	"fmt"
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
	// 1. Configurações mínimas obrigatórias
	hd.ConfigGravarValor("Principal", "TipoResposta", "2") // 2 = JSON
	hd.ConfigGravarValor("Principal", "LogNivel", "3")
	hd.ConfigGravarValor("Principal", "LogPath", "/app/logs/acbr/")

	hd.ConfigGravarValor("DFe", "SSLCryptLib", "1") // OpenSSL
	hd.ConfigGravarValor("DFe", "SSLHttpLib", "3")
	hd.ConfigGravarValor("DFe", "SSLXmlSignLib", "4")
	hd.ConfigGravarValor("DFe", "UF", "BA")
	hd.ConfigGravarValor("DFe", "ArquivoPFX", "/app/certs/empresa.pfx")
	hd.ConfigGravarValor("DFe", "Senha", "suaSenha")
	hd.ConfigGravarValor("DFe", "Timeout", "15000")

	hd.ConfigGravarValor("NFe", "ModeloDF", "55")
	hd.ConfigGravarValor("NFe", "VersaoDF", "4.00")
	hd.ConfigGravarValor("NFe", "Ambiente", "2") // 2=Homologação

	hd.ConfigGravarValor("Arquivos", "PathSalvar", "/app/data/nfe/")
	hd.ConfigGravarValor("Arquivos", "PathSchemas", "/app/lib/Schemas/NFe/")
	hd.ConfigGravarValor("Arquivos", "Salvar", "1")

	// 3. Se precisar persistir depois:
	// hd.ConfigGravar("/tmp/acbr_final.ini")

	hd.ConfiguredFor = companyID
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
