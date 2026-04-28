package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InvoiceItem represents a product line in an invoice.
type InvoiceItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	InvoiceID uuid.UUID `json:"invoice_id" db:"invoice_id"`
	NItem     int       `json:"n_item" db:"n_item"`

	// Produto
	CProd    string           `json:"c_prod" db:"c_prod"`
	CEAN     string           `json:"c_ean" db:"c_ean"`
	XProd    string           `json:"x_prod" db:"x_prod"`
	NCM      string           `json:"ncm" db:"ncm"`
	CEST     string           `json:"cest,omitempty" db:"cest"`
	CFOP     string           `json:"cfop" db:"cfop"`
	UCom     string           `json:"u_com" db:"u_com"`
	QCom     decimal.Decimal  `json:"q_com" db:"q_com"`
	VUnCom   decimal.Decimal  `json:"v_un_com" db:"v_un_com"`
	VProd    decimal.Decimal  `json:"v_prod" db:"v_prod"`
	CEANTrib string           `json:"c_ean_trib,omitempty" db:"c_ean_trib"`
	UTrib    string           `json:"u_trib,omitempty" db:"u_trib"`
	QTrib    *decimal.Decimal `json:"q_trib,omitempty" db:"q_trib"`
	VUnTrib  *decimal.Decimal `json:"v_un_trib,omitempty" db:"v_un_trib"`
	VFrete   decimal.Decimal  `json:"v_frete" db:"v_frete"`
	VSeg     decimal.Decimal  `json:"v_seg" db:"v_seg"`
	VDesc    decimal.Decimal  `json:"v_desc" db:"v_desc"`
	VOutro   decimal.Decimal  `json:"v_outro" db:"v_outro"`
	IndTot   int16            `json:"ind_tot" db:"ind_tot"`
	XPed     string           `json:"x_ped,omitempty" db:"x_ped"`
	NItemPed string           `json:"n_item_ped,omitempty" db:"n_item_ped"`
	VTotTrib *decimal.Decimal `json:"v_tot_trib,omitempty" db:"v_tot_trib"`
	InfAdProd string          `json:"inf_ad_prod,omitempty" db:"inf_ad_prod"`
	IndEscala string          `json:"ind_escala,omitempty" db:"ind_escala"`
	CNPJFab   string          `json:"cnpj_fab,omitempty" db:"cnpj_fab"`
	CBenef    string          `json:"c_benef,omitempty" db:"c_benef"`

	// Reforma - Produto
	IndBemMovelUsado *int16           `json:"ind_bem_movel_usado,omitempty" db:"ind_bem_movel_usado"`
	VItem            *decimal.Decimal `json:"v_item,omitempty" db:"v_item"`

	// ICMS (campos principais - todos nullable exceto orig)
	ICMSOrig     *int16  `json:"icms_orig,omitempty" db:"icms_orig"`
	ICMSCST      string  `json:"icms_cst,omitempty" db:"icms_cst"`
	ICMSCSOSN    string  `json:"icms_csosn,omitempty" db:"icms_csosn"`
	ICMSModBC    *int16  `json:"icms_mod_bc,omitempty" db:"icms_mod_bc"`
	ICMSPRedBC   *decimal.Decimal `json:"icms_p_red_bc,omitempty" db:"icms_p_red_bc"`
	ICMSVBC      *decimal.Decimal `json:"icms_v_bc,omitempty" db:"icms_v_bc"`
	ICMSPICMS    *decimal.Decimal `json:"icms_p_icms,omitempty" db:"icms_p_icms"`
	ICMSVICMS    *decimal.Decimal `json:"icms_v_icms,omitempty" db:"icms_v_icms"`
	ICMSModBCST  *int16           `json:"icms_mod_bc_st,omitempty" db:"icms_mod_bc_st"`
	ICMSPMVAST   *decimal.Decimal `json:"icms_p_mva_st,omitempty" db:"icms_p_mva_st"`
	ICMSPRedBCST *decimal.Decimal `json:"icms_p_red_bc_st,omitempty" db:"icms_p_red_bc_st"`
	ICMSVBCST    *decimal.Decimal `json:"icms_v_bc_st,omitempty" db:"icms_v_bc_st"`
	ICMSPICMSST  *decimal.Decimal `json:"icms_p_icms_st,omitempty" db:"icms_p_icms_st"`
	ICMSVICMSST  *decimal.Decimal `json:"icms_v_icms_st,omitempty" db:"icms_v_icms_st"`
	ICMSUFST     string           `json:"icms_uf_st,omitempty" db:"icms_uf_st"`
	ICMSPBCOp    *decimal.Decimal `json:"icms_p_bc_op,omitempty" db:"icms_p_bc_op"`
	ICMSVBCSTRet *decimal.Decimal `json:"icms_v_bc_st_ret,omitempty" db:"icms_v_bc_st_ret"`
	ICMSVICMSSTRet *decimal.Decimal `json:"icms_v_icms_st_ret,omitempty" db:"icms_v_icms_st_ret"`
	ICMSMotDes     *int16           `json:"icms_mot_des,omitempty" db:"icms_mot_des"`
	ICMSPCredSN    *decimal.Decimal `json:"icms_p_cred_sn,omitempty" db:"icms_p_cred_sn"`
	ICMSVCredICMSSN *decimal.Decimal `json:"icms_v_cred_icms_sn,omitempty" db:"icms_v_cred_icms_sn"`
	ICMSVICMSDeson  *decimal.Decimal `json:"icms_v_icms_deson,omitempty" db:"icms_v_icms_deson"`
	ICMSVICMSOp     *decimal.Decimal `json:"icms_v_icms_op,omitempty" db:"icms_v_icms_op"`
	ICMSPDif        *decimal.Decimal `json:"icms_p_dif,omitempty" db:"icms_p_dif"`
	ICMSVICMSDif    *decimal.Decimal `json:"icms_v_icms_dif,omitempty" db:"icms_v_icms_dif"`
	ICMSPST         *decimal.Decimal `json:"icms_p_st,omitempty" db:"icms_p_st"`
	ICMSVBCFCP      *decimal.Decimal `json:"icms_v_bc_fcp,omitempty" db:"icms_v_bc_fcp"`
	ICMSPFCP        *decimal.Decimal `json:"icms_p_fcp,omitempty" db:"icms_p_fcp"`
	ICMSVFCP        *decimal.Decimal `json:"icms_v_fcp,omitempty" db:"icms_v_fcp"`
	ICMSVBCFCPSTRet *decimal.Decimal `json:"icms_v_bc_fcp_st_ret,omitempty" db:"icms_v_bc_fcp_st_ret"`
	ICMSPFCPSTRet   *decimal.Decimal `json:"icms_p_fcp_st_ret,omitempty" db:"icms_p_fcp_st_ret"`
	ICMSVFCPSTRet   *decimal.Decimal `json:"icms_v_fcp_st_ret,omitempty" db:"icms_v_fcp_st_ret"`
	ICMSPRedBCEfet  *decimal.Decimal `json:"icms_p_red_bc_efet,omitempty" db:"icms_p_red_bc_efet"`
	ICMSVBCEfet     *decimal.Decimal `json:"icms_v_bc_efet,omitempty" db:"icms_v_bc_efet"`
	ICMSPICMSEfet   *decimal.Decimal `json:"icms_p_icms_efet,omitempty" db:"icms_p_icms_efet"`
	ICMSVICMSEfet   *decimal.Decimal `json:"icms_v_icms_efet,omitempty" db:"icms_v_icms_efet"`
	ICMSVICMSSub    *decimal.Decimal `json:"icms_v_icms_substituto,omitempty" db:"icms_v_icms_substituto"`

	// ICMS Monofásico
	ICMSQBCMono        *decimal.Decimal `json:"icms_q_bc_mono,omitempty" db:"icms_q_bc_mono"`
	ICMSAdRemICMS      *decimal.Decimal `json:"icms_ad_rem_icms,omitempty" db:"icms_ad_rem_icms"`
	ICMSVICMSMono      *decimal.Decimal `json:"icms_v_icms_mono,omitempty" db:"icms_v_icms_mono"`
	ICMSQBCMonoReten   *decimal.Decimal `json:"icms_q_bc_mono_reten,omitempty" db:"icms_q_bc_mono_reten"`
	ICMSAdRemICMSReten *decimal.Decimal `json:"icms_ad_rem_icms_reten,omitempty" db:"icms_ad_rem_icms_reten"`
	ICMSVICMSMonoReten *decimal.Decimal `json:"icms_v_icms_mono_reten,omitempty" db:"icms_v_icms_mono_reten"`
	ICMSPRedAdRem      *decimal.Decimal `json:"icms_p_red_ad_rem,omitempty" db:"icms_p_red_ad_rem"`
	ICMSMotRedAdRem    *int16           `json:"icms_mot_red_ad_rem,omitempty" db:"icms_mot_red_ad_rem"`
	ICMSQBCMonoRet     *decimal.Decimal `json:"icms_q_bc_mono_ret,omitempty" db:"icms_q_bc_mono_ret"`
	ICMSVICMSMonoOp    *decimal.Decimal `json:"icms_v_icms_mono_op,omitempty" db:"icms_v_icms_mono_op"`
	ICMSVICMSMonoDif   *decimal.Decimal `json:"icms_v_icms_mono_dif,omitempty" db:"icms_v_icms_mono_dif"`
	ICMSAdRemICMSRet   *decimal.Decimal `json:"icms_ad_rem_icms_ret,omitempty" db:"icms_ad_rem_icms_ret"`
	ICMSVICMSMonoRet   *decimal.Decimal `json:"icms_v_icms_mono_ret,omitempty" db:"icms_v_icms_mono_ret"`

	// ICMS UF Dest
	ICMSUFDestVBC         *decimal.Decimal `json:"icms_uf_dest_v_bc,omitempty" db:"icms_uf_dest_v_bc"`
	ICMSUFDestVBCFCP      *decimal.Decimal `json:"icms_uf_dest_v_bc_fcp,omitempty" db:"icms_uf_dest_v_bc_fcp"`
	ICMSUFDestPFCP        *decimal.Decimal `json:"icms_uf_dest_p_fcp,omitempty" db:"icms_uf_dest_p_fcp"`
	ICMSUFDestPICMS       *decimal.Decimal `json:"icms_uf_dest_p_icms,omitempty" db:"icms_uf_dest_p_icms"`
	ICMSUFDestPICMSInter  *decimal.Decimal `json:"icms_uf_dest_p_icms_inter,omitempty" db:"icms_uf_dest_p_icms_inter"`
	ICMSUFDestPICMSInterP *decimal.Decimal `json:"icms_uf_dest_p_icms_inter_part,omitempty" db:"icms_uf_dest_p_icms_inter_part"`
	ICMSUFDestVFCP        *decimal.Decimal `json:"icms_uf_dest_v_fcp,omitempty" db:"icms_uf_dest_v_fcp"`
	ICMSUFDestVICMS       *decimal.Decimal `json:"icms_uf_dest_v_icms,omitempty" db:"icms_uf_dest_v_icms"`
	ICMSUFDestVICMSRemet  *decimal.Decimal `json:"icms_uf_dest_v_icms_remet,omitempty" db:"icms_uf_dest_v_icms_remet"`

	// PIS
	PISCST      string           `json:"pis_cst,omitempty" db:"pis_cst"`
	PISVBC      *decimal.Decimal `json:"pis_v_bc,omitempty" db:"pis_v_bc"`
	PISPPIS     *decimal.Decimal `json:"pis_p_pis,omitempty" db:"pis_p_pis"`
	PISQBCProd  *decimal.Decimal `json:"pis_q_bc_prod,omitempty" db:"pis_q_bc_prod"`
	PISVAliqProd *decimal.Decimal `json:"pis_v_aliq_prod,omitempty" db:"pis_v_aliq_prod"`
	PISVPIS     *decimal.Decimal `json:"pis_v_pis,omitempty" db:"pis_v_pis"`

	// COFINS
	COFINSCST      string           `json:"cofins_cst,omitempty" db:"cofins_cst"`
	COFINSVBC      *decimal.Decimal `json:"cofins_v_bc,omitempty" db:"cofins_v_bc"`
	COFINSPCOFINS  *decimal.Decimal `json:"cofins_p_cofins,omitempty" db:"cofins_p_cofins"`
	COFINSQBCProd  *decimal.Decimal `json:"cofins_q_bc_prod,omitempty" db:"cofins_q_bc_prod"`
	COFINSVAliqProd *decimal.Decimal `json:"cofins_v_aliq_prod,omitempty" db:"cofins_v_aliq_prod"`
	COFINSVCONFINS *decimal.Decimal `json:"cofins_v_cofins,omitempty" db:"cofins_v_cofins"`

	// II
	IIVBC     *decimal.Decimal `json:"ii_v_bc,omitempty" db:"ii_v_bc"`
	IIVDespAdu *decimal.Decimal `json:"ii_v_desp_adu,omitempty" db:"ii_v_desp_adu"`
	IIVII     *decimal.Decimal `json:"ii_v_ii,omitempty" db:"ii_v_ii"`
	IIVIOF    *decimal.Decimal `json:"ii_v_iof,omitempty" db:"ii_v_iof"`

	// IPI
	IPICST  string           `json:"ipi_cst,omitempty" db:"ipi_cst"`
	IPICEnq string           `json:"ipi_c_enq,omitempty" db:"ipi_c_enq"`
	IPIVBC  *decimal.Decimal `json:"ipi_v_bc,omitempty" db:"ipi_v_bc"`
	IPIPIPI *decimal.Decimal `json:"ipi_p_ipi,omitempty" db:"ipi_p_ipi"`
	IPIVIPI *decimal.Decimal `json:"ipi_v_ipi,omitempty" db:"ipi_v_ipi"`

	// IS (Reforma)
	ISCST        string           `json:"is_cst,omitempty" db:"is_cst"`
	ISCClassTrib string           `json:"is_c_class_trib,omitempty" db:"is_c_class_trib"`
	ISVBC        *decimal.Decimal `json:"is_v_bc,omitempty" db:"is_v_bc"`
	ISPIS        *decimal.Decimal `json:"is_p_is,omitempty" db:"is_p_is"`
	ISPISEspec   *decimal.Decimal `json:"is_p_is_espec,omitempty" db:"is_p_is_espec"`
	ISVIS        *decimal.Decimal `json:"is_v_is,omitempty" db:"is_v_is"`

	// IBS/CBS (Reforma)
	IBSCBSCST        string           `json:"ibs_cbs_cst,omitempty" db:"ibs_cbs_cst"`
	IBSCBSCClassTrib string           `json:"ibs_cbs_c_class_trib,omitempty" db:"ibs_cbs_c_class_trib"`
	IBSCBSVBC        *decimal.Decimal `json:"ibs_cbs_v_bc,omitempty" db:"ibs_cbs_v_bc"`
	IBSCBSVIBS       *decimal.Decimal `json:"ibs_cbs_v_ibs,omitempty" db:"ibs_cbs_v_ibs"`
	IBSUFP           *decimal.Decimal `json:"ibs_uf_p,omitempty" db:"ibs_uf_p"`
	IBSUFV           *decimal.Decimal `json:"ibs_uf_v,omitempty" db:"ibs_uf_v"`
	IBSUFPDif        *decimal.Decimal `json:"ibs_uf_p_dif,omitempty" db:"ibs_uf_p_dif"`
	IBSUFVDif        *decimal.Decimal `json:"ibs_uf_v_dif,omitempty" db:"ibs_uf_v_dif"`
	IBSUFVDevTrib    *decimal.Decimal `json:"ibs_uf_v_dev_trib,omitempty" db:"ibs_uf_v_dev_trib"`
	IBSUFPRedAliq    *decimal.Decimal `json:"ibs_uf_p_red_aliq,omitempty" db:"ibs_uf_p_red_aliq"`
	IBSUFPAliqEfet   *decimal.Decimal `json:"ibs_uf_p_aliq_efet,omitempty" db:"ibs_uf_p_aliq_efet"`
	IBSMunP          *decimal.Decimal `json:"ibs_mun_p,omitempty" db:"ibs_mun_p"`
	IBSMunV          *decimal.Decimal `json:"ibs_mun_v,omitempty" db:"ibs_mun_v"`
	IBSMunPDif       *decimal.Decimal `json:"ibs_mun_p_dif,omitempty" db:"ibs_mun_p_dif"`
	IBSMunVDif       *decimal.Decimal `json:"ibs_mun_v_dif,omitempty" db:"ibs_mun_v_dif"`
	IBSMunVDevTrib   *decimal.Decimal `json:"ibs_mun_v_dev_trib,omitempty" db:"ibs_mun_v_dev_trib"`
	IBSMunPRedAliq   *decimal.Decimal `json:"ibs_mun_p_red_aliq,omitempty" db:"ibs_mun_p_red_aliq"`
	IBSMunPAliqEfet  *decimal.Decimal `json:"ibs_mun_p_aliq_efet,omitempty" db:"ibs_mun_p_aliq_efet"`
	CBSP             *decimal.Decimal `json:"cbs_p,omitempty" db:"cbs_p"`
	CBSV             *decimal.Decimal `json:"cbs_v,omitempty" db:"cbs_v"`
	CBSPDif          *decimal.Decimal `json:"cbs_p_dif,omitempty" db:"cbs_p_dif"`
	CBSVDif          *decimal.Decimal `json:"cbs_v_dif,omitempty" db:"cbs_v_dif"`
	CBSVDevTrib      *decimal.Decimal `json:"cbs_v_dev_trib,omitempty" db:"cbs_v_dev_trib"`
	CBSPRedAliq      *decimal.Decimal `json:"cbs_p_red_aliq,omitempty" db:"cbs_p_red_aliq"`
	CBSPAliqEfet     *decimal.Decimal `json:"cbs_p_aliq_efet,omitempty" db:"cbs_p_aliq_efet"`

	// Tributação Regular (Reforma)
	TribRegCST          string           `json:"trib_reg_cst,omitempty" db:"trib_reg_cst"`
	TribRegCClassTrib   string           `json:"trib_reg_c_class_trib,omitempty" db:"trib_reg_c_class_trib"`
	TribRegPAliqIBSUF   *decimal.Decimal `json:"trib_reg_p_aliq_ibs_uf,omitempty" db:"trib_reg_p_aliq_ibs_uf"`
	TribRegVIBSUF       *decimal.Decimal `json:"trib_reg_v_ibs_uf,omitempty" db:"trib_reg_v_ibs_uf"`
	TribRegPAliqIBSMun  *decimal.Decimal `json:"trib_reg_p_aliq_ibs_mun,omitempty" db:"trib_reg_p_aliq_ibs_mun"`
	TribRegVIBSMun      *decimal.Decimal `json:"trib_reg_v_ibs_mun,omitempty" db:"trib_reg_v_ibs_mun"`
	TribRegPAliqCBS     *decimal.Decimal `json:"trib_reg_p_aliq_cbs,omitempty" db:"trib_reg_p_aliq_cbs"`
	TribRegVCBS         *decimal.Decimal `json:"trib_reg_v_cbs,omitempty" db:"trib_reg_v_cbs"`

	// Crédito Presumido (Reforma)
	IBSCredPresCod string           `json:"ibs_cred_pres_cod,omitempty" db:"ibs_cred_pres_cod"`
	IBSCredPresP   *decimal.Decimal `json:"ibs_cred_pres_p,omitempty" db:"ibs_cred_pres_p"`
	IBSCredPresV   *decimal.Decimal `json:"ibs_cred_pres_v,omitempty" db:"ibs_cred_pres_v"`
	CBSCredPresCod string           `json:"cbs_cred_pres_cod,omitempty" db:"cbs_cred_pres_cod"`
	CBSCredPresP   *decimal.Decimal `json:"cbs_cred_pres_p,omitempty" db:"cbs_cred_pres_p"`
	CBSCredPresV   *decimal.Decimal `json:"cbs_cred_pres_v,omitempty" db:"cbs_cred_pres_v"`

	// Compra Governamental (Reforma)
	GovPAliqIBSUF  *decimal.Decimal `json:"gov_p_aliq_ibs_uf,omitempty" db:"gov_p_aliq_ibs_uf"`
	GovVTribIBSUF  *decimal.Decimal `json:"gov_v_trib_ibs_uf,omitempty" db:"gov_v_trib_ibs_uf"`
	GovPAliqIBSMun *decimal.Decimal `json:"gov_p_aliq_ibs_mun,omitempty" db:"gov_p_aliq_ibs_mun"`
	GovVTribIBSMun *decimal.Decimal `json:"gov_v_trib_ibs_mun,omitempty" db:"gov_v_trib_ibs_mun"`
	GovPAliqCBS    *decimal.Decimal `json:"gov_p_aliq_cbs,omitempty" db:"gov_p_aliq_cbs"`
	GovVTribCBS    *decimal.Decimal `json:"gov_v_trib_cbs,omitempty" db:"gov_v_trib_cbs"`

	// Monofásico IBS/CBS (Reforma)
	MonoVTotIBS       *decimal.Decimal `json:"mono_v_tot_ibs,omitempty" db:"mono_v_tot_ibs"`
	MonoVTotCBS       *decimal.Decimal `json:"mono_v_tot_cbs,omitempty" db:"mono_v_tot_cbs"`
	MonoPadQBC        *decimal.Decimal `json:"mono_pad_q_bc,omitempty" db:"mono_pad_q_bc"`
	MonoPadAdRemIBS   *decimal.Decimal `json:"mono_pad_ad_rem_ibs,omitempty" db:"mono_pad_ad_rem_ibs"`
	MonoPadAdRemCBS   *decimal.Decimal `json:"mono_pad_ad_rem_cbs,omitempty" db:"mono_pad_ad_rem_cbs"`
	MonoPadVIBS       *decimal.Decimal `json:"mono_pad_v_ibs,omitempty" db:"mono_pad_v_ibs"`
	MonoPadVCBS       *decimal.Decimal `json:"mono_pad_v_cbs,omitempty" db:"mono_pad_v_cbs"`
	MonoRetenQBC      *decimal.Decimal `json:"mono_reten_q_bc,omitempty" db:"mono_reten_q_bc"`
	MonoRetenAdRemIBS *decimal.Decimal `json:"mono_reten_ad_rem_ibs,omitempty" db:"mono_reten_ad_rem_ibs"`
	MonoRetenVIBS     *decimal.Decimal `json:"mono_reten_v_ibs,omitempty" db:"mono_reten_v_ibs"`
	MonoRetenAdRemCBS *decimal.Decimal `json:"mono_reten_ad_rem_cbs,omitempty" db:"mono_reten_ad_rem_cbs"`
	MonoRetenVCBS     *decimal.Decimal `json:"mono_reten_v_cbs,omitempty" db:"mono_reten_v_cbs"`
	MonoRetQBC        *decimal.Decimal `json:"mono_ret_q_bc,omitempty" db:"mono_ret_q_bc"`
	MonoRetAdRemIBS   *decimal.Decimal `json:"mono_ret_ad_rem_ibs,omitempty" db:"mono_ret_ad_rem_ibs"`
	MonoRetVIBS       *decimal.Decimal `json:"mono_ret_v_ibs,omitempty" db:"mono_ret_v_ibs"`
	MonoRetAdRemCBS   *decimal.Decimal `json:"mono_ret_ad_rem_cbs,omitempty" db:"mono_ret_ad_rem_cbs"`
	MonoRetVCBS       *decimal.Decimal `json:"mono_ret_v_cbs,omitempty" db:"mono_ret_v_cbs"`
	MonoDifPIBS       *decimal.Decimal `json:"mono_dif_p_ibs,omitempty" db:"mono_dif_p_ibs"`
	MonoDifVIBS       *decimal.Decimal `json:"mono_dif_v_ibs,omitempty" db:"mono_dif_v_ibs"`
	MonoDifPCBS       *decimal.Decimal `json:"mono_dif_p_cbs,omitempty" db:"mono_dif_p_cbs"`
	MonoDifVCBS       *decimal.Decimal `json:"mono_dif_v_cbs,omitempty" db:"mono_dif_v_cbs"`

	// Transferência Crédito / ZFM (Reforma)
	TransfCredVIBS      *decimal.Decimal `json:"transf_cred_v_ibs,omitempty" db:"transf_cred_v_ibs"`
	TransfCredVCBS      *decimal.Decimal `json:"transf_cred_v_cbs,omitempty" db:"transf_cred_v_cbs"`
	CredPresIBSZFMTp    string           `json:"cred_pres_ibs_zfm_tp,omitempty" db:"cred_pres_ibs_zfm_tp"`
	CredPresIBSZFMV     *decimal.Decimal `json:"cred_pres_ibs_zfm_v,omitempty" db:"cred_pres_ibs_zfm_v"`
}
