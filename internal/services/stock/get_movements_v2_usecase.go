package services

import (
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *stockAuditService) GetAllMovements(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error) {
	return s.stockRepo.GetAllMovements(filter, companyID)
}
