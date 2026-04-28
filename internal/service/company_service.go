package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/validator"
)

type CompanyService interface {
	Create(ctx context.Context, userID uuid.UUID, req *dto.CompanyCreateRequest) (*dto.CompanyResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.CompanyResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.CompanyCreateRequest) (*dto.CompanyResponse, error)
	ConfigureSMTP(ctx context.Context, id uuid.UUID, req *dto.CompanySMTPRequest) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]dto.CompanyResponse, error)
	LinkUser(ctx context.Context, ownerID, userID, companyID uuid.UUID) error
}

type companyService struct {
	repo       repository.CompanyRepository
	certRepo   repository.CertificateRepository
	cryptoSvc  crypto.Service
}

func NewCompanyService(repo repository.CompanyRepository, certRepo repository.CertificateRepository, cryptoSvc crypto.Service) CompanyService {
	return &companyService{
		repo:      repo,
		certRepo:  certRepo,
		cryptoSvc: cryptoSvc,
	}
}

func (s *companyService) Create(ctx context.Context, userID uuid.UUID, req *dto.CompanyCreateRequest) (*dto.CompanyResponse, error) {
	if !validator.IsValidCNPJ(req.CNPJ) {
		return nil, apperror.NewBadRequest("invalid CNPJ")
	}

	company := &domain.Company{
		ID:                 uuid.New(),
		CNPJ:               req.CNPJ,
		RazaoSocial:        req.RazaoSocial,
		NomeFantasia:       req.NomeFantasia,
		InscricaoEstadual:  req.InscricaoEstadual,
		InscricaoMunicipal: req.InscricaoMunicipal,
		CRT:                req.CRT,
		Logradouro:         req.Logradouro,
		Numero:             req.Numero,
		Complemento:        req.Complemento,
		Bairro:             req.Bairro,
		CodMunicipio:       req.CodMunicipio,
		Municipio:          req.Municipio,
		UF:                 req.UF,
		CEP:                req.CEP,
		Telefone:           req.Telefone,
		CNAE:               req.CNAE,
		Ambiente:           req.Ambiente,
		SerieNFe:           req.SerieNFe,
		SerieNFCe:          req.SerieNFCe,
		CSCID:              req.CSCID,
		CSCToken:           req.CSCToken,
		IsActive:           true,
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, err
	}

	if err := s.repo.LinkUser(ctx, userID, company.ID); err != nil {
		return nil, err
	}

	return s.mapToResponse(company, false), nil
}

func (s *companyService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CompanyResponse, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, errCert := s.certRepo.GetByCompanyID(ctx, id)
	hasCert := errCert == nil

	return s.mapToResponse(company, hasCert), nil
}

func (s *companyService) Update(ctx context.Context, id uuid.UUID, req *dto.CompanyCreateRequest) (*dto.CompanyResponse, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields (CNPJ cannot be updated)
	company.RazaoSocial = req.RazaoSocial
	company.NomeFantasia = req.NomeFantasia
	company.InscricaoEstadual = req.InscricaoEstadual
	company.InscricaoMunicipal = req.InscricaoMunicipal
	company.CRT = req.CRT
	company.Logradouro = req.Logradouro
	company.Numero = req.Numero
	company.Complemento = req.Complemento
	company.Bairro = req.Bairro
	company.CodMunicipio = req.CodMunicipio
	company.Municipio = req.Municipio
	company.UF = req.UF
	company.CEP = req.CEP
	company.Telefone = req.Telefone
	company.CNAE = req.CNAE
	company.Ambiente = req.Ambiente
	company.SerieNFe = req.SerieNFe
	company.SerieNFCe = req.SerieNFCe
	company.CSCID = req.CSCID
	company.CSCToken = req.CSCToken

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, err
	}

	_, errCert := s.certRepo.GetByCompanyID(ctx, id)
	hasCert := errCert == nil

	return s.mapToResponse(company, hasCert), nil
}

func (s *companyService) ConfigureSMTP(ctx context.Context, id uuid.UUID, req *dto.CompanySMTPRequest) error {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	encPassword, err := s.cryptoSvc.Encrypt([]byte(req.SMTPPassword))
	if err != nil {
		return err
	}

	company.SMTPHost = req.SMTPHost
	company.SMTPPort = req.SMTPPort
	company.SMTPUser = req.SMTPUser
	company.SMTPPasswordEnc = encPassword
	company.SMTPFrom = req.SMTPFrom
	company.SMTPTLS = *req.SMTPTLS

	return s.repo.Update(ctx, company)
}

func (s *companyService) ListByUser(ctx context.Context, userID uuid.UUID) ([]dto.CompanyResponse, error) {
	companies, err := s.repo.GetCompaniesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []dto.CompanyResponse
	for _, c := range companies {
		_, errCert := s.certRepo.GetByCompanyID(ctx, c.ID)
		res = append(res, *s.mapToResponse(&c, errCert == nil))
	}
	return res, nil
}

func (s *companyService) LinkUser(ctx context.Context, ownerID, userID, companyID uuid.UUID) error {
	// Need to check if ownerID actually owns companyID (this can be checked in middleware as well)
	return s.repo.LinkUser(ctx, userID, companyID)
}

func (s *companyService) mapToResponse(c *domain.Company, hasCert bool) *dto.CompanyResponse {
	return &dto.CompanyResponse{
		ID:                 c.ID,
		CNPJ:               c.CNPJ,
		RazaoSocial:        c.RazaoSocial,
		NomeFantasia:       c.NomeFantasia,
		InscricaoEstadual:  c.InscricaoEstadual,
		InscricaoMunicipal: c.InscricaoMunicipal,
		CRT:                c.CRT,
		Logradouro:         c.Logradouro,
		Numero:             c.Numero,
		Complemento:        c.Complemento,
		Bairro:             c.Bairro,
		CodMunicipio:       c.CodMunicipio,
		Municipio:          c.Municipio,
		UF:                 c.UF,
		CEP:                c.CEP,
		Telefone:           c.Telefone,
		CNAE:               c.CNAE,
		Ambiente:           c.Ambiente,
		SerieNFe:           c.SerieNFe,
		SerieNFCe:          c.SerieNFCe,
		CSCID:              c.CSCID,
		SMTPConfigured:     c.SMTPHost != "",
		HasCertificate:     hasCert,
		IsActive:           c.IsActive,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
