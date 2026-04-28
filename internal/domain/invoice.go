package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Invoice represents an emitted NFe/NFCe header.
type Invoice struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CompanyID uuid.UUID `json:"company_id" db:"company_id"`

	// --- Chave de acesso ---
	Chave string `json:"chave" db:"chave"`

	// --- Identificação ---
	CNF        string    `json:"c_nf" db:"c_nf"`
	NatOp      string    `json:"nat_op" db:"nat_op"`
	Modelo     int16     `json:"modelo" db:"modelo"`
	Serie      int       `json:"serie" db:"serie"`
	Numero     int       `json:"numero" db:"numero"`
	DHEmissao  time.Time `json:"dh_emissao" db:"dh_emissao"`
	DHSaiEnt   *time.Time `json:"dh_sai_ent,omitempty" db:"dh_sai_ent"`
	TpNF       int16     `json:"tp_nf" db:"tp_nf"`
	IDDest     int16     `json:"id_dest" db:"id_dest"`
	CMunFG     string    `json:"c_mun_fg" db:"c_mun_fg"`
	TpImp      int16     `json:"tp_imp" db:"tp_imp"`
	TpEmis     int16     `json:"tp_emis" db:"tp_emis"`
	TpAmb      int16     `json:"tp_amb" db:"tp_amb"`
	FinNFe     int16     `json:"fin_nfe" db:"fin_nfe"`
	IndFinal   int16     `json:"ind_final" db:"ind_final"`
	IndPres    int16     `json:"ind_pres" db:"ind_pres"`
	ProcEmi    int16     `json:"proc_emi" db:"proc_emi"`
	VerProc    string    `json:"ver_proc" db:"ver_proc"`
	IndIntermed string   `json:"ind_intermed,omitempty" db:"ind_intermed"`

	// --- Reforma Tributária (Identificação) ---
	CMunFGIBS   string           `json:"c_mun_fg_ibs,omitempty" db:"c_mun_fg_ibs"`
	TpNFDebito  string           `json:"tp_nf_debito,omitempty" db:"tp_nf_debito"`
	TpNFCredito string           `json:"tp_nf_credito,omitempty" db:"tp_nf_credito"`
	TpEnteGov   *int16           `json:"tp_ente_gov,omitempty" db:"tp_ente_gov"`
	PRedutor    *decimal.Decimal `json:"p_redutor,omitempty" db:"p_redutor"`
	TpOperGov   *int16           `json:"tp_oper_gov,omitempty" db:"tp_oper_gov"`

	// --- Destinatário ---
	DestCNPJCPF      string `json:"dest_cnpj_cpf,omitempty" db:"dest_cnpj_cpf"`
	DestNome         string `json:"dest_nome,omitempty" db:"dest_nome"`
	DestIE           string `json:"dest_ie,omitempty" db:"dest_ie"`
	DestIndIEDest    *int16 `json:"dest_ind_ie_dest,omitempty" db:"dest_ind_ie_dest"`
	DestEmail        string `json:"dest_email,omitempty" db:"dest_email"`
	DestLogradouro   string `json:"dest_logradouro,omitempty" db:"dest_logradouro"`
	DestNumero       string `json:"dest_numero,omitempty" db:"dest_numero"`
	DestComplemento  string `json:"dest_complemento,omitempty" db:"dest_complemento"`
	DestBairro       string `json:"dest_bairro,omitempty" db:"dest_bairro"`
	DestCodMunicipio string `json:"dest_cod_municipio,omitempty" db:"dest_cod_municipio"`
	DestMunicipio    string `json:"dest_municipio,omitempty" db:"dest_municipio"`
	DestUF           string `json:"dest_uf,omitempty" db:"dest_uf"`
	DestCEP          string `json:"dest_cep,omitempty" db:"dest_cep"`

	// --- Totais ICMS ---
	TotVBC         decimal.Decimal `json:"tot_v_bc" db:"tot_v_bc"`
	TotVICMS       decimal.Decimal `json:"tot_v_icms" db:"tot_v_icms"`
	TotVICMSDeson  decimal.Decimal `json:"tot_v_icms_deson" db:"tot_v_icms_deson"`
	TotVFCP        decimal.Decimal `json:"tot_v_fcp" db:"tot_v_fcp"`
	TotVBCST       decimal.Decimal `json:"tot_v_bc_st" db:"tot_v_bc_st"`
	TotVST         decimal.Decimal `json:"tot_v_st" db:"tot_v_st"`
	TotVFCPST      decimal.Decimal `json:"tot_v_fcp_st" db:"tot_v_fcp_st"`
	TotVProd       decimal.Decimal `json:"tot_v_prod" db:"tot_v_prod"`
	TotVFrete      decimal.Decimal `json:"tot_v_frete" db:"tot_v_frete"`
	TotVSeg        decimal.Decimal `json:"tot_v_seg" db:"tot_v_seg"`
	TotVDesc       decimal.Decimal `json:"tot_v_desc" db:"tot_v_desc"`
	TotVII         decimal.Decimal `json:"tot_v_ii" db:"tot_v_ii"`
	TotVIPI        decimal.Decimal `json:"tot_v_ipi" db:"tot_v_ipi"`
	TotVPIS        decimal.Decimal `json:"tot_v_pis" db:"tot_v_pis"`
	TotVCOFINS     decimal.Decimal `json:"tot_v_cofins" db:"tot_v_cofins"`
	TotVOutro      decimal.Decimal `json:"tot_v_outro" db:"tot_v_outro"`
	TotVNF         decimal.Decimal `json:"tot_v_nf" db:"tot_v_nf"`
	TotVTotTrib    decimal.Decimal `json:"tot_v_tot_trib" db:"tot_v_tot_trib"`
	TotVIPIDevol   decimal.Decimal `json:"tot_v_ipi_devol" db:"tot_v_ipi_devol"`
	TotVFCPSTRet   decimal.Decimal `json:"tot_v_fcp_st_ret" db:"tot_v_fcp_st_ret"`

	// --- Totais Monofásico ICMS ---
	TotQBCMono         *decimal.Decimal `json:"tot_q_bc_mono,omitempty" db:"tot_q_bc_mono"`
	TotVICMSMono       *decimal.Decimal `json:"tot_v_icms_mono,omitempty" db:"tot_v_icms_mono"`
	TotQBCMonoReten    *decimal.Decimal `json:"tot_q_bc_mono_reten,omitempty" db:"tot_q_bc_mono_reten"`
	TotVICMSMonoReten  *decimal.Decimal `json:"tot_v_icms_mono_reten,omitempty" db:"tot_v_icms_mono_reten"`
	TotQBCMonoRet      *decimal.Decimal `json:"tot_q_bc_mono_ret,omitempty" db:"tot_q_bc_mono_ret"`
	TotVICMSMonoRet    *decimal.Decimal `json:"tot_v_icms_mono_ret,omitempty" db:"tot_v_icms_mono_ret"`

	// --- Totais Reforma Tributária ---
	TotVNFTot                *decimal.Decimal `json:"tot_v_nf_tot,omitempty" db:"tot_v_nf_tot"`
	TotISVIS                 *decimal.Decimal `json:"tot_is_v_is,omitempty" db:"tot_is_v_is"`
	TotIBSCBSVBC             *decimal.Decimal `json:"tot_ibs_cbs_v_bc,omitempty" db:"tot_ibs_cbs_v_bc"`
	TotIBSVIBS               *decimal.Decimal `json:"tot_ibs_v_ibs,omitempty" db:"tot_ibs_v_ibs"`
	TotIBSVCredPres          *decimal.Decimal `json:"tot_ibs_v_cred_pres,omitempty" db:"tot_ibs_v_cred_pres"`
	TotIBSVCredPresCondSus   *decimal.Decimal `json:"tot_ibs_v_cred_pres_cond_sus,omitempty" db:"tot_ibs_v_cred_pres_cond_sus"`
	TotIBSUFVDif             *decimal.Decimal `json:"tot_ibs_uf_v_dif,omitempty" db:"tot_ibs_uf_v_dif"`
	TotIBSUFVDevTrib         *decimal.Decimal `json:"tot_ibs_uf_v_dev_trib,omitempty" db:"tot_ibs_uf_v_dev_trib"`
	TotIBSUFVIBSUF           *decimal.Decimal `json:"tot_ibs_uf_v_ibs_uf,omitempty" db:"tot_ibs_uf_v_ibs_uf"`
	TotIBSMunVDif            *decimal.Decimal `json:"tot_ibs_mun_v_dif,omitempty" db:"tot_ibs_mun_v_dif"`
	TotIBSMunVDevTrib        *decimal.Decimal `json:"tot_ibs_mun_v_dev_trib,omitempty" db:"tot_ibs_mun_v_dev_trib"`
	TotIBSMunVIBSMun         *decimal.Decimal `json:"tot_ibs_mun_v_ibs_mun,omitempty" db:"tot_ibs_mun_v_ibs_mun"`
	TotCBSVDif               *decimal.Decimal `json:"tot_cbs_v_dif,omitempty" db:"tot_cbs_v_dif"`
	TotCBSVDevTrib           *decimal.Decimal `json:"tot_cbs_v_dev_trib,omitempty" db:"tot_cbs_v_dev_trib"`
	TotCBSVCBS               *decimal.Decimal `json:"tot_cbs_v_cbs,omitempty" db:"tot_cbs_v_cbs"`
	TotCBSVCredPres          *decimal.Decimal `json:"tot_cbs_v_cred_pres,omitempty" db:"tot_cbs_v_cred_pres"`
	TotCBSVCredPresCondSus   *decimal.Decimal `json:"tot_cbs_v_cred_pres_cond_sus,omitempty" db:"tot_cbs_v_cred_pres_cond_sus"`
	TotMonoVIBSMono          *decimal.Decimal `json:"tot_mono_v_ibs_mono,omitempty" db:"tot_mono_v_ibs_mono"`
	TotMonoVCBSMono          *decimal.Decimal `json:"tot_mono_v_cbs_mono,omitempty" db:"tot_mono_v_cbs_mono"`
	TotMonoVIBSMonoReten     *decimal.Decimal `json:"tot_mono_v_ibs_mono_reten,omitempty" db:"tot_mono_v_ibs_mono_reten"`
	TotMonoVCBSMonoReten     *decimal.Decimal `json:"tot_mono_v_cbs_mono_reten,omitempty" db:"tot_mono_v_cbs_mono_reten"`
	TotMonoVIBSMonoRet       *decimal.Decimal `json:"tot_mono_v_ibs_mono_ret,omitempty" db:"tot_mono_v_ibs_mono_ret"`
	TotMonoVCBSMonoRet       *decimal.Decimal `json:"tot_mono_v_cbs_mono_ret,omitempty" db:"tot_mono_v_cbs_mono_ret"`

	// --- Protocolo SEFAZ ---
	Protocolo      string     `json:"protocolo,omitempty" db:"protocolo"`
	DHRecebimento  *time.Time `json:"dh_recebimento,omitempty" db:"dh_recebimento"`
	Status         string     `json:"status" db:"status"` // autorizada, cancelada, denegada, inutilizada

	// --- Informações Adicionais ---
	InfAdFisco string `json:"inf_ad_fisco,omitempty" db:"inf_ad_fisco"`
	InfCpl     string `json:"inf_cpl,omitempty" db:"inf_cpl"`

	// --- B2 Storage ---
	XMLB2Key string `json:"xml_b2_key,omitempty" db:"xml_b2_key"`
	PDFB2Key string `json:"pdf_b2_key,omitempty" db:"pdf_b2_key"`

	// --- Timestamps ---
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// --- Relations (populated by queries) ---
	Items     []InvoiceItem     `json:"items,omitempty" db:"-"`
	Payments  []InvoicePayment  `json:"payments,omitempty" db:"-"`
	Transport *InvoiceTransport `json:"transport,omitempty" db:"-"`
}
