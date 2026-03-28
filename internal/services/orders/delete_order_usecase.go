package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func DeleteOrderUseCase(id int, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()

	// Use a transaction to ensure consistency
	return db.Transaction(func(tx *gorm.DB) error {
		var order types.ProductionOrder

		if err := tx.Preload("Items").Where("id_company = ? AND id_order = ?", companyID, id).First(&order).Error; err != nil {
			return fmt.Errorf("order not found: %w", err)
		}

		for _, item := range order.Items {
			if item.MachineID != nil && *item.MachineID != 0 {
				if err := tx.Model(&types.Machine{}).Where("id = ?", *item.MachineID).Update("status", "idle").Error; err != nil {
					return fmt.Errorf("failed to set machine to idle: %w", err)
				}
				services.PublishEvent("dashboard_updates", `{"type": "MACHINE_STATUS_CHANGED", "message": "Machine set to idle due to order deletion"}`)
			}
		}

		if err := tx.Delete(&order).Error; err != nil {
			return fmt.Errorf("failed to delete order: %w", err)
		}

		services.InvalidateCache("orders:list:*")
		services.PublishEvent("dashboard_updates", `{"type": "ORDER_DELETED", "message": "ORDER DELETED"}`)

		return nil
	})
}
