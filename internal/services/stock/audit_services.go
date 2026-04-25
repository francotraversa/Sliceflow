package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
)

type stockAuditService struct {
	stockRepo domain.StockAuditRepository
}

func NewStockAuditService(stockRepo domain.StockAuditRepository) domain.StockAuditService {
	return &stockAuditService{stockRepo: stockRepo}
}

func (s *stockAuditService) GetMovementByID(id uint, companyID uint) (*types.StockMovement, error) {
	return s.stockRepo.GetMovementByID(id, companyID)
}
