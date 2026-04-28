package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
)

type InvoiceService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceListResponse, error) // In a real app we would map the full domain.Invoice
}

type invoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) GetByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceListResponse, error) {
	inv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.InvoiceListResponse{
		ID:        inv.ID,
		Chave:     inv.Chave,
		Numero:    inv.Numero,
		Serie:     inv.Serie,
		Modelo:    inv.Modelo,
		DHEmissao: inv.DHEmissao,
		Status:    inv.Status,
		TotVNF:    inv.TotVNF,
		DestNome:  inv.DestNome,
		DestCNPJCPF: inv.DestCNPJCPF,
		CreatedAt: inv.CreatedAt,
	}, nil
}
