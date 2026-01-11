package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func UpdateOrderUseCase(id int, dto types.UpdateOrderDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder

	// 1. Cargamos la orden con sus ítems actuales para que GORM conozca el estado previo
	if err := db.Preload("Items").First(&order, id).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 2. Actualizamos los campos de la tabla 'production_orders'
	if dto.ClientName != nil {
		order.ClientName = *dto.ClientName
	}
	if dto.Priority != nil {
		order.Priority = *dto.Priority
	}
	if dto.Notes != nil {
		order.Notes = *dto.Notes
	}
	if dto.Status != nil {
		order.Status = *dto.Status
	}
	if dto.Price != nil {
		order.Price = dto.Price
	}
	if dto.OperatorID != nil {
		order.OperatorID = *dto.OperatorID
	}
	if dto.MaterialID != nil {
		order.MaterialID = *dto.MaterialID
	}
	if dto.MachineID != nil {
		order.MachineID = dto.MachineID
	}

	// 3. Sincronizamos la tabla 'order_items'
	if dto.Items != nil {
		var updatedItems []types.OrderItem
		currentTotalDone := 0
		currentTotalPieces := 0

		for _, itemDTO := range *dto.Items {

			item := types.OrderItem{
				ID:         itemDTO.ID,
				OrderID:    order.ID,
				StlName:    itemDTO.ProductName,
				Quantity:   itemDTO.Quantity,
				DonePieces: itemDTO.DonePieces,
			}
			updatedItems = append(updatedItems, item)
			currentTotalDone += item.DonePieces
			currentTotalPieces += item.Quantity
		}
		order.Items = updatedItems
		order.DonePieces = currentTotalDone
		order.TotalPieces = currentTotalPieces
	}

	// 4. Lógica de tiempos y auto-completado
	if order.DonePieces >= order.TotalPieces && order.TotalPieces > 0 {
		order.Status = "completed"
	}

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error

	if err != nil {
		return fmt.Errorf("failed to update order and items: %w", err)
	}

	// 6. Limpieza de cache y eventos
	services.InvalidateCache("orders:list:*")
	services.PublishEvent("dashboard_updates", `{"type": "ORDER_UPDATED", "message": "ORDER UPDATED"}`)

	return nil
}
