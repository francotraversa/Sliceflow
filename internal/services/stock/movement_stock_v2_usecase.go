package services

import (
	"errors"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *stockAuditService) CreateMovement(req types.CreateMovementRequest, companyID uint) error {
	if req.SKU == "" {
		return errors.New("SKU is required")
	}
	if req.Type == "" {
		return errors.New("Type is required")
	}
	if req.Quantity == 0 {
		return errors.New("Quantity is required")
	}
	if req.UserID == 0 {
		return errors.New("User ID is required")
	}
	if req.LocationID == 0 {
		return errors.New("Location ID is required")
	}
	if req.Description == "" {
		return errors.New("Description is required")
	}
	if req.Reason == "" {
		return errors.New("Reason is required")
	}
	return s.stockRepo.ExecuteMovementTransaction(req, companyID)
}
