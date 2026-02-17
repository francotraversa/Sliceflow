package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func DeleteOrderUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder
	if err := db.First(&order, id).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}
	if err := db.Delete(&order).Error; err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	services.InvalidateCache("orders:list:*")
	services.PublishEvent("dashboard_updates", `{"type": "ORDER_DELETED", "message": "ORDER DELETED"}`)
	return nil
}
