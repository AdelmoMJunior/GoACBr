# Guia Completo ACBrLibNFe – 100% Métodos, Configuração e Exemplos

**Versão do documento:** abril/2026  
**Baseado na documentação oficial ACBrLibNFe 4.00**

---

## 1. Métodos da Biblioteca (46 funções exportadas)

### 1.1 Inicialização e Informações
| Método | Comando | Descrição |
|---|---|---|
| Inicializar | `NFE_Inicializar(eArqConfig, eChaveCrypt)` | Sempre chamar primeiro. Cria INI se vazio【6811202811139311301†L31-L35】 |
| Finalizar | `NFE_Finalizar()` | Libera memória【2827649645496921650†L7-L10】 |
| Nome | `NFE_Nome(sResposta, esTamanho)` | Retorna "ACBrLibNFe" |
| Versão | `NFE_Versao(sResposta, esTamanho)` | Retorna versão da DLL |
| ÚltimoRetorno | `NFE_UltimoRetorno(sResposta, esTamanho)` | Texto do último erro/sucesso |

### 1.2 Configuração
| Método | Comando |
|---|---|
| ConfigLer | `NFE_ConfigLer(eArqConfig, sResposta, esTamanho)` |
| ConfigGravar | `NFE_ConfigGravar(eArqConfig, sResposta, esTamanho)` |
| ConfigLerValor | `NFE_ConfigLerValor(eSessao, eChave, sResposta, esTamanho)` |
| ConfigGravarValor | `NFE_ConfigGravarValor(eSessao, eChave, sValor)`【3779636941541182480†L7-L10】 |
| ConfigImportar | `NFE_ConfigImportar(eArquivoIni)` |
| ConfigExportar | `NFE_ConfigExportar(sResposta, esTamanho)` |

### 1.3 Carga e Manipulação
| Método | Comando |
|---|---|
| CarregarXML | `NFE_CarregarXML(eArquivoOuXML)` |
| CarregarINI | `NFE_CarregarINI(eArquivoOuINI)` |
| CarregarEventoXML | `NFE_CarregarEventoXML(eArquivoOuXML)` |
| CarregarEventoINI | `NFE_CarregarEventoINI(eArquivoOuINI)` |
| LimparLista | `NFE_LimparLista()` |
| LimparListaEventos | `NFE_LimparListaEventos()` |
| Assinar | `NFE_Assinar()` |
| Validar | `NFE_Validar()` |
| ValidarRegrasdeNegocios | `NFE_ValidarRegrasdeNegocios(sResposta, esTamanho)` |
| VerificarAssinatura | `NFE_VerificarAssinatura(sResposta, esTamanho)` |
| ObterXml | `NFE_ObterXml(AIndex, sResposta, esTamanho)`【1591023077566579260†L20-L21】 |
| GravarXml | `NFE_GravarXml(AIndex, eNomeArquivo, ePathArquivo)`【1591023077566579260†L21-L22】 |
| ObterIni | `NFE_ObterIni(AIndex, sResposta, esTamanho)`【1591023077566579260†L22-L23】 |
| GravarIni | `NFE_GravarIni(AIndex, eNomeArquivo, ePathArquivo)` |

### 1.4 WebServices
| Método | Comando |
|---|---|
| StatusServico | `NFE_StatusServico(sResposta, esTamanho)` |
| Consultar | `NFE_Consultar(eChaveOuNFe, sResposta, esTamanho)` |
| Inutilizar | `NFE_Inutilizar(ACNPJ, AJustificativa, Ano, Modelo, Serie, NumeroInicial, NumeroFinal, sResposta, esTamanho)` |
| Enviar | `NFE_Enviar(aLote, bImprimir, bSincrono, bZipado, sResposta, esTamanho)` |
| Cancelar | `NFE_Cancelar(eChave, eJustificativa, eCNPJ, ALote, sResposta, esTamanho)` |
| EnviarEvento | `NFE_EnviarEvento(idLote, sResposta, esTamanho)` |

### 1.5 Distribuição DFe
| Método | Comando |
|---|---|
| DistribuicaoDFe | `NFE_DistribuicaoDFe(AcUFAutor, ACNPJCPF, nUltNSU, nNSU, sResposta, esTamanho)` |
| DistribuicaoDFePorUltNSU | `NFE_DistribuicaoDFePorUltNSU(AcUFAutor, ACNPJCPF, eUltNSU, sResposta, esTamanho)` |
| DistribuicaoDFePorNSU | `NFE_DistribuicaoDFePorNSU(AcUFAutor, ACNPJCPF, eNSU, sResposta, esTamanho)` |
| DistribuicaoDFePorChave | `NFE_DistribuicaoDFePorChave(AcUFAutor, ACNPJCPF, eChave, sResposta, esTamanho)` |

### 1.6 Email
| Método | Comando |
|---|---|
| EnviarEmail | `NFE_EnviarEmail(ePara, eXmlNFe, bEnviaPDF, eAssunto, eCC, eAnexos, eMensagem, sResposta, esTamanho)` |
| EnviarEmailEvento | `NFE_EnviarEmailEvento(ePara, eXmlEvento, eXmlNFe, bEnviaPDF, eAssunto, eCC, eAnexos, eMensagem, sResposta, esTamanho)` |

### 1.7 Impressão
| Método | Comando |
|---|---|
| Imprimir | `NFE_Imprimir(cImpressora, nNumCopias, cProtocolo, bMostrarPreview, cMarcaDagua, bViaConsumidor, bSimplificado)`【793792565539354369†L8-L10】 |
| ImprimirPDF | `NFE_ImprimirPDF()` |
| ImprimirEvento | `NFE_ImprimirEvento(eArquivoXmlNFe, eArquivoXmlEvento)` |
| ImprimirEventoPDF | `NFE_ImprimirEventoPDF(eArquivoXmlNFe, eArquivoXmlEvento)`【3912040403234196676†L8-L10】 |
| ImprimirInutilizacao | `NFE_ImprimirInutilizacao(eArquivoXml)` |
| ImprimirInutilizacaoPDF | `NFE_ImprimirInutilizacaoPDF(eArquivoXml)` |
| SalvarPDF | `NFE_SalvarPDF(sResposta, esTamanho)`【3912040403234196676†L27-L29】 |

### 1.8 Utilitários
| Método | Comando |
|---|---|
| ObterCertificados | `NFE_ObterCertificados(sResposta, esTamanho)`【6811202811139311301†L70-L73】 |

---

## 2. Configuração Completa ACBrLib.ini

### [Principal]
TipoResposta=2 ; 0=INI 1=XML 2=JSON
CodificacaoResposta=0 ; 0=UTF8
LogNivel=3 ; 0-4
LogPath=C:\ACBr\Log\

### [DFe]
SSLLib=1 ; 1=OpenSSL
CryptLib=1
HttpLib=3
XmlSignLib=4
SSLType=5 ; TLS 1.2
ArquivoPFX=C:\certs\empresa.pfx
Senha=123456
NumeroSerie=
VerificarValidade=1
AguardarConsultaRet=0
IntervaloTentativas=1000
Tentativas=5
Timeout=15000
QuebradeLinha=|

### [NFe]
FormaEmissao=0 ; 0=Normal
ModeloDF=55 ; 55 ou 65
VersaoDF=4.00
IdCSC=
CSC=
SalvarTXT=0
SalvarXML=1
SalvarEvento=1
SalvarApenasNFeProcessadas=1
EmissaoPathNFe=1
NormatizarMunicipios=1
PathNFe=
PathInu=
PathEvento=

### [WebService]
UF=SP
Ambiente=2 ; 1=Producao 2=Homologacao
Visualizar=0
Salvar=1
AjustaAguarda=1
Aguardar=5000
Tentativas=5
IntervaloTentativas=1000

### [Proxy]
Servidor=
Porta=
Usuario=
Senha=

### [Arquivos]
Salvar=1
SepararPorMes=1
SepararPorCNPJ=1
SepararPorModelo=1
AdicionarLiteral=1
EmissaoPathNFe=1
PathSalvar=C:\ACBr\NFe\
PathSchemas=C:\ACBr\Schemas\NFe\

### [Email]
Nome=Empresa Ltda
Conta=nfe@empresa.com.br
Servidor=smtp.empresa.com.br
Porta=587
Usuario=nfe@empresa.com.br
Senha=senha
SSL=0
TLS=1
Autenticar=1

### [SoftwareHouse]
CNPJ=11222333000144
Nome=Minha Software House
Fone=11999999999
Email=contato@software.com.br

### [PosPrinter]
Modelo=1
Porta=USB
PaginaDeCodigo=0
Colunas=48
EspacoEntreLinhas=0

---

## 3. Modelo INI NFe 55 – Completo

```ini
[infNFe]
versao=4.00

[Identificacao]
cNF=12345678
natOp=Venda de Mercadoria
indPag=0
mod=55
serie=1
nNF=1001
dhEmi=28/04/2026 09:30:00
dhSaiEnt=28/04/2026 09:30:00
tpNF=1
idDest=1
cMunFG=3550308
tpImp=1
tpEmis=1
tpAmb=2
finNFe=1
indFinal=1
indPres=1
procEmi=0
verProc=ACBrNFe

[Emitente]
CRT=1
CNPJCPF=11222333000144
xNome=EMPRESA MODELO LTDA
xFant=EMPRESA MODELO
IE=123456789012
IM=12345
CNAE=4711302
xLgr=Rua Cel Aureliano Camargo
nro=973
xBairro=Centro
cMun=3554003
xMun=Tatui
UF=SP
CEP=18270000
cPais=1058
xPais=BRASIL
Fone=11999999999

[Destinatario]
CNPJCPF=99999999000199
xNome=CLIENTE TESTE
indIEDest=1
IE=987654321098
Email=cliente@teste.com.br
xLgr=Rua das Flores
nro=123
xBairro=Centro
cMun=3550308
xMun=Sao Paulo
UF=SP
CEP=04615000

[Produto001]
cProd=0001
cEAN=SEM GTIN
xProd=NOTEBOOK I5
NCM=84713012
CFOP=5102
uCom=UN
qCom=1
vUnCom=3500.00
vProd=3500.00
uTrib=UN
qTrib=1
vUnTrib=3500.00
indTot=1

[ICMS001]
orig=0
CSOSN=102
pCredSN=0
vCredICMSSN=0

[PIS001]
CST=07

[COFINS001]
CST=07

[Total]
vProd=3500.00
vNF=3500.00
vTotTrib=500.00

[Transportador]
modFrete=1

[pag001]
tPag=01
vPag=3500.00
indPag=0

[infAdic]
infCpl=Documento emitido por ME optante Simples Nacional
```

---

## 4. Modelo INI NFCe 65 – Completo

```ini
[infNFe]
versao=4.00

[Identificacao]
cNF=87654321
natOp=Venda
mod=65
serie=1
nNF=5200
dhEmi=28/04/2026 10:15:00
tpNF=1
idDest=1
cMunFG=3550308
tpImp=4
tpEmis=1
tpAmb=2
finNFe=1
indFinal=1
indPres=1

[Emitente]
CRT=1
CNPJCPF=11222333000144
xNome=LOJA MODELO
IE=123456789012
xLgr=Av Paulista
nro=1000
xBairro=Bela Vista
cMun=3550308
xMun=Sao Paulo
UF=SP
CEP=01310000

[Produto001]
cProd=2001
cEAN=7891234567890
xProd=REFRIGERANTE 2L
NCM=22021000
CFOP=5405
uCom=UN
qCom=2
vUnCom=8.50
vProd=17.00
indTot=1

[ICMS001]
orig=0
CSOSN=102

[Total]
vProd=17.00
vNF=17.00

[pag001]
tPag=03
vPag=17.00
tBand=01
tpIntegra=1
```

---

## 5. Eventos – INI para cada tipo

### Cancelamento (110111)
```ini
[EVENTO]
idLote=1
[EVENTO001]
cOrgao=35
CNPJ=11222333000144
chNFe=35260411222333000144550010000010011000000001
dhEvento=2026-04-28T10:15:00-03:00
tpEvento=110111
nSeqEvento=1
versaoEvento=1.00
nProt=135240000000000
xJust=Erro no valor do produto
```

### Carta Correção (110110)
```ini
[EVENTO001]
cOrgao=35
CNPJ=11222333000144
chNFe=...
dhEvento=2026-04-28T10:20:00-03:00
tpEvento=110110
nSeqEvento=1
versaoEvento=1.00
xCorrecao=O CFOP correto e 5102
```

### Manifestação – Ciência (210210)
```ini
[EVENTO001]
cOrgao=91
CNPJ=99999999000199
chNFe=...
dhEvento=2026-04-28T10:30:00-03:00
tpEvento=210210
nSeqEvento=1
versaoEvento=1.00
```

### EPEC (110140)
```ini
[EVENTO001]
cOrgao=35
CNPJ=11222333000144
chNFe=...
dhEvento=2026-04-28T10:40:00-03:00
tpEvento=110140
nSeqEvento=1
versaoEvento=1.00
cOrgaoAutor=35
tpAutor=1
verAplic=ACBrLibNFe
dhEmi=2026-04-28T10:39:00-03:00
tpNF=1
IE=123456789012
destUF=SP
vNF=1500.00
vICMS=180.00
vST=0.00
[DEST]
DestCNPJCPF=99999999000199
```

---

## 6. Fluxo típico em código

```pascal
NFE_Inicializar('C:\ACBr\ACBrLib.ini','');
NFE_ConfigGravarValor('DFe','ArquivoPFX','C:\certs\empresa.pfx');
NFE_ConfigGravarValor('DFe','Senha','123456');
NFE_CarregarINI('c:\temp\nfe.ini');
NFE_Assinar();
NFE_Validar();
NFE_Enviar(1, False, True, False, resposta, tamanho);
NFE_ImprimirPDF();
NFE_Finalizar();
```

Todos os métodos, configurações e modelos acima foram extraídos da documentação oficial ACBrLibNFe.
