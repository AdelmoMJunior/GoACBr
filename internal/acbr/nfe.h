#ifndef NFE_H
#define NFE_H

#include <stdint.h>

/* Definindo tipos para compatibilidade com o Go/CGO */
typedef void* handle;

/* --- Inicialização e Finalização --- */
int NFE_Inicializar(handle* h, const char* eArqConfig, const char* eChaveCrypt);
int NFE_Finalizar(handle h);

/* --- Configurações --- */
int NFE_ConfigLer(handle h, const char* eArqConfig);
int NFE_ConfigGravar(handle h, const char* eArqConfig);
int NFE_ConfigLerValor(handle h, const char* eSessao, const char* eChave, char* sValor, int* esTamanho);
int NFE_ConfigGravarValor(handle h, const char* eSessao, const char* eChave, const char* sValor);

/* --- Operações Principais NFe/NFCe --- */
int NFE_CarregarINI(handle h, const char* eArquivoOuString);
int NFE_LimparLista(handle h);
int NFE_Assinar(handle h);
int NFE_Validar(handle h);
int NFE_Enviar(handle h, int aLote, int imprimir, int sincrono, int zipado, char* sResposta, int* esTamanho);

/* --- Consultas e Eventos --- */
int NFE_Consultar(handle h, const char* eChaveOuNFe, int aExtrairEventos, char* sResposta, int* esTamanho);
int NFE_Cancelar(handle h, const char* eChave, const char* eJustificativa, const char* eCNPJ, int aLote, char* sResposta, int* esTamanho);
int NFE_CartaCorrecao(handle h, const char* eChave, const char* eCorrecao, const char* eCNPJ, int aLote, char* sResposta, int* esTamanho);
int NFE_Inutilizar(handle h, const char* eCNPJ, const char* eJustificativa, int aAno, int aModelo, int aSerie, int aNumeroInicial, int aNumeroFinal, char* sResposta, int* esTamanho);
int NFE_DistribuicaoDFePorUltNSU(handle h, int acUFAutor, const char* eCNPJCPF, const char* eultNSU, char* sResposta, int* esTamanho);
int NFE_DistribuicaoDFePorNSU(handle h, int acUFAutor, const char* eCNPJCPF, const char* eNSU, char* sResposta, int* esTamanho);

/* --- PDFs e XMLs --- */
int NFE_ImprimirPDF(handle h);
int NFE_ObterXml(handle h, int aIndex, char* sResposta, int* esTamanho);
int NFE_ObterCaminhoGerado(handle h, char* sResposta, int* esTamanho);

/* --- Status e Último Retorno --- */
int NFE_UltimoRetorno(handle h, char* sMensagem, int* esTamanho);

#endif

