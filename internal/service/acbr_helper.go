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

	slog.Debug("Applying company configuration via INI file", "company_id", companyID)

	// Ensure valid values for numeric fields
	amb := comp.Ambiente
	if amb != 1 && amb != 2 {
		amb = 2 // Default to Homologação
	}

	// Build INI Content - EXACTLY matching user's template but with Linux paths
	iniContent := fmt.Sprintf(`[Principal]
TipoResposta=0
CodificacaoResposta=0
LogNivel=0
LogPath=/app/logs/acbr

[Versao]
ACBrLib=0.0.2
ACBrLibNFE=1.4.7.415

[Sistema]
Nome=AM SOFTware
Versao=0.0.1
Data=30/12/1899
Descricao=

[Email]
Nome=
Servidor=
Conta=
Usuario=
Senha=
Codificacao=27
Porta=0
SSL=0
TLS=0
SSLType=5
Timeout=0
Confirmacao=0
ConfirmacaoEntrega=0
SegundoPlano=0
Tentativas=1
IsHTML=0
Priority=1

[PosPrinter]
ArqLog=
Modelo=0
Porta=
PaginaDeCodigo=2
ColunasFonteNormal=48
EspacoEntreLinhas=0
LinhasEntreCupons=21
CortaPapel=1
TraduzirTags=1
IgnorarTags=0
LinhasBuffer=0
ControlePorta=0
VerificarImpressora=0
TipoCorte=0

[PosPrinter_Barras]
MostrarCodigo=0
LarguraLinha=0
Altura=0
Margem=0

[PosPrinter_QRCode]
Tipo=2
LarguraModulo=4
ErrorLevel=0

[PosPrinter_Logo]
IgnorarLogo=0
KeyCode1=32
KeyCode2=32
FatorX=1
FatorY=1

[PosPrinter_Gaveta]
SinalInvertido=0
TempoON=50
TempoOFF=200

[PosPrinter_MPagina]
Largura=0
Altura=0
Esquerda=0
Topo=0
Direcao=0
EspacoEntreLinhas=0

[PosPrinter_Device]
Baud=9600
Data=8
Timeout=3
Parity=0
Stop=0
MaxBandwidth=0
SendBytesCount=0
SendBytesInterval=0
HandShake=0
SoftFlow=0
HardFlow=0

[Proxy]
Servidor=
Porta=0
Usuario=
Senha=

[Socket]
NivelLog=0
ArqLog=
Timeout=0

[SoftwareHouse]
CNPJ=
RazaoSocial=
NomeFantasia=
WebSite=
Email=
Telefone=
Responsavel=

[Emissor]
CNPJ=
RazaoSocial=
NomeFantasia=
WebSite=
Email=
Telefone=
Responsavel=

[DFe]
SSLCryptLib=1
SSLHttpLib=3
SSLXmlSignLib=4
UF=%s
TimeZone.Modo=0
TimeZone.Str=
URLPFX=
ArquivoPFX=%s
DadosPFX=
Senha=%s
NumeroSerie=
VerificarValidade=1

[NFe]
FormaEmissao=0
SalvarGer=1
ExibirErroSchema=1
FormatoAlerta=TAG:%%TAGNIVEL%% ID:%%ID%%/%%TAG%%(%%DESCRICAO%%) - %%MSG%%.
RetirarAcentos=1
RetirarEspacos=1
IdentarXML=0
ValidarDigest=1
IdCSC=
CSC=
ModeloDF=55
VersaoDF=4.00
AtualizarXMLCancelado=1
VersaoQRCode=2
CamposFatObrigatorios=1
TagNT2018005=0
ForcarGerarTagRejeicao906=0
Ambiente=%d
SalvarWS=1
Timeout=15000
TimeoutPorThread=0
Visualizar=0
AjustaAguardaConsultaRet=1
AguardarConsultaRet=5000
IntervaloTentativas=1000
Tentativas=5
SSLType=5
QuebradeLinha=|
PathSalvar=/app/data/nfe
PathSchemas=/app/lib/Schemas/NFe
IniServicos=
SalvarArq=1
AdicionarLiteral=0
SepararPorCNPJ=1
SepararPorIE=0
SepararPorModelo=1
SepararPorAno=1
SepararPorMes=1
SepararPorDia=0
Download.PathDownload=/app/data/nfe/download
Download.SepararPorNome=0
SalvarEvento=1
SalvarApenasNFeProcessadas=1
EmissaoPathNFe=0
NormatizarMunicipios=0
PathNFe=/app/data/nfe/xml
PathInu=/app/data/nfe/inutilizacao
PathEvento=/app/data/nfe/evento
PathArquivoMunicipios=/app/data/nfe/municipio
IdCSRT=0
CSRT=

[WebService]
UF=%s
Ambiente=%d
Visualizar=0
Salvar=1
AjustaAguarda=1
Aguardar=5000
Tentativas=5
IntervaloTentativas=3000
TimeZone=-3

[Arquivos]
Salvar=1
SepararPorMes=1
SepararPorCNPJ=1
SepararPorModelo=1
AdicionarLiteral=1
PathSalvar=/app/data/nfe
PathSchemas=/app/lib/Schemas/NFe

[DANFE]
PathPDF=/app/data/nfe/pdf
UsaSeparadorPathPDF=0
Impressora=
NomeDocumento=
MostraSetup=0
MostraPreview=1
MostraStatus=1
Copias=1
PathLogo=
MargemInferior=8
MargemSuperior=8
MargemEsquerda=6
MargemDireita=5
AlterarEscalaPadrao=0
NovaEscala=96
ExpandeLogoMarca=0
ExpandeLogoMarca.Altura=0
ExpandeLogoMarca.Esquerda=0
ExpandeLogoMarca.Topo=0
ExpandeLogoMarca.Largura=0
ExpandeLogoMarca.Dimensionar=0
ExpandeLogoMarca.Esticar=1
CasasDecimais.Formato=0
CasasDecimais.MaskqCom=,0.00
CasasDecimais.MaskvUnCom=,0.00
CasasDecimais.qCom=2
CasasDecimais.vUnCom=2
CasasDecimais.MaskAliquota=,0.00
CasasDecimais.Aliquota=2
Protocolo=
Cancelada=0
TipoDANFE=1
ImprimeTotalLiquido=1
vTribFed=0
vTribEst=0
vTribMun=0
FonteTributos=
ChaveTributos=
ImprimeTributos=1
ExibeTotalTributosItem=0
ImprimeCodigoEan=0
ImprimeNomeFantasia=0
ExibeInforAdicProduto=1
QuebraLinhaEmDetalhamentos=1

[DANFENFe]
FormularioContinuo=0
ImprimeValor=0
ImprimeDescPorPercentual=0
ImprimeDetalhamentoEspecifico=1
PosCanhoto=0
PosCanhotoLayout=0
ExibeResumoCanhoto=1
TextoResumoCanhoto=
ExibeCampoFatura=1
ExibeDadosISSQN=0
ExibeDadosDocReferenciados=1
DetVeiculos=[dv_chassi,dv_xCor,dv_nSerie,dv_tpComb,dv_nMotor,dv_anoMod,dv_anoFab]
DetMedicamentos=[dm_nLote,dm_qLote,dm_dFab,dm_dVal,dm_vPMC]
DetArmamentos=[da_tpArma,da_nSerie,da_nCano,da_descr]
DetCombustiveis=[dc_cProdANP,dc_CODIF,dc_qTemp,dc_UFCons,dc_CIDE,dc_qBCProd,dc_vAliqProd,dc_vCIDE]
DetRastros=[dr_nLote,dr_qLote,dr_dFab,dr_dVal,dr_cAgreg]
TributosPercentual=0
TributosPercentualPersonalizado=0
MarcadAgua=
LarguraCodProd=54
ExibeEAN=0
AltLinhaComun=30
EspacoEntreProdutos=7
AlternaCoresProdutos=0
TamanhoLogoHeight=0
TamanhoLogoWidth=0
RecuoEndereco=0
RecuoEmpresa=0
LogoemCima=0
RecuoLogo=0
ExpandirDadosAdicionaisAuto=0
ImprimeContDadosAdPrimeiraPagina=0
ExibeCampoDePagamento=0
ImprimeInscSuframa=1
ImprimeXPedNitemPed=0
ImprimeDescAcrescItemNFe=0
FormatarNumeroDocumento=1
CorDestaqueProdutos=clWhite
Fonte.Nome=0
Fonte.Negrito=0
Fonte.TamanhoFonteRazaoSocial=8
Fonte.TamanhoFonteEndereco=0
Fonte.TamanhoFonteInformacoesComplementares=8
Fonte.TamanhoFonteDemaisCampos=8

[DANFENFCe]
TipoRelatorioBobina=0
TipoRelatorioEvento=0
LarguraBobina=302
ImprimeDescAcrescItem=1
ImprimeItens=1
ViaConsumidor=0
vTroco=0
ImprimeQRCodeLateral=0
ImprimeLogoLateral=0
EspacoFinal=38
TamanhoLogoHeight=50
TamanhoLogoWidth=77
DescricaoPagamentos=[icaTipo,icaBandeira]
ImprimeEmUmaLinha=0
ImprimeEmDuasLinhas=0
MargemInferior=0
MargemSuperior=0
MargemEsquerda=0
MargemDireita=0
FonteLinhaItem.Name=Lucida Console
FonteLinhaItem.Color=536870912
FonteLinhaItem.Size=7
FonteLinhaItem.Bold=0
FonteLinhaItem.Italic=0
FonteLinhaItem.Underline=0
FonteLinhaItem.StrikeOut=0
FormatarNumeroDocumento=1
`, comp.UF, pfxPath, password, amb, comp.UF, amb)

	slog.Debug("Generated INI content (Turbo Mode)", "ini", iniContent)

	// ACBrLib is picky about line endings even on Linux, use \r\n
	iniContent = strings.ReplaceAll(iniContent, "\n", "\r\n")

	// Ensure directories exist before loading config
	_ = os.MkdirAll("/app/logs/acbr", 0755)
	_ = os.MkdirAll("/app/data/nfe", 0755)

	// Write to a temporary company-specific INI file
	tmpIniPath := fmt.Sprintf("/tmp/acbr_%s.ini", companyID)
	if err := os.WriteFile(tmpIniPath, []byte(iniContent), 0666); err != nil {
		return fmt.Errorf("failed to write temporary ini: %w", err)
	}
	defer os.Remove(tmpIniPath)

	slog.Debug("Temporary INI written successfully", "path", tmpIniPath)

	// List schemas directory to verify deployment
	schemaPath := "/app/lib/Schemas/NFe"
	entries, err := os.ReadDir(schemaPath)
	if err != nil {
		slog.Error("Failed to read schemas directory", "path", schemaPath, "error", err)
		// Try listing parent to see what happened
		parentEntries, _ := os.ReadDir("/app/lib/Schemas")
		var pFiles []string
		for _, e := range parentEntries {
			pFiles = append(pFiles, e.Name())
		}
		slog.Debug("Files in parent Schemas directory", "files", pFiles)
	} else {
		var files []string
		for _, e := range entries {
			files = append(files, e.Name())
		}
		slog.Debug("Files in Schemas directory", "path", schemaPath, "count", len(files), "files", files)
	}

	if err := hd.ConfigLer(tmpIniPath); err != nil {
		return err
	}

	hd.ConfiguredFor = companyID
	return nil
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
