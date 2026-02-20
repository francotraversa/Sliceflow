package services

import (
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	servicesMachine "github.com/francotraversa/Sliceflow/internal/services/machine"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func UpdateOrderUseCase(id int, dto types.UpdateOrderDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder
	var oldMachineID *int
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
		oldMachineID = order.MachineID
		order.MachineID = dto.MachineID
	}

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

	if order.Status != "completed" {
		// Auto-status basado en piezas
		if order.DonePieces >= order.TotalPieces {
			order.Status = "ready"
		} else if order.DonePieces > 0 {
			order.Status = "in-progress"
		}

		if dto.MachineID != nil && (oldMachineID == nil || *oldMachineID != *dto.MachineID) {
			if oldMachineID != nil {
				db.Model(&types.Machine{}).Where("id = ?", *oldMachineID).Update("status", "idle")
			}
			if *dto.MachineID != 0 {
				db.Model(&types.Machine{}).Where("id = ?", *dto.MachineID).Update("status", "printing")
			}
		}
	} else {
		if order.MachineID != nil {
			newStatus := "idle"
			updmachine := types.UpdateMachineDTO{Status: &newStatus}
			servicesMachine.UpdateMachineUseCase(*order.MachineID, updmachine)
		}
		now := time.Now()
		order.FinishTime = &now
	}
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error

	if err != nil {
		return fmt.Errorf("failed to update order and items: %w", err)
	}

	servicesWeb.InvalidateCache("orders:list:*")
	servicesWeb.PublishEvent("dashboard_updates", `{"type": "ORDER_UPDATED", "message": "ORDER UPDATED"}`)

	return nil
}
