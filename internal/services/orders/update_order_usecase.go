package services

import (
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	servicesMachine "github.com/francotraversa/Sliceflow/internal/services/machine"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateOrderUseCase(id int, dto types.UpdateOrderDTO, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder

	// Buscar por id (PK de la BD) con filtro de compañía para seguridad multitenant
	if err := db.Preload("Items").Where("id_company = ? AND id = ?", companyID, id).First(&order).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Actualizar campos del order
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
	if dto.OperatorID != nil {
		order.OperatorID = *dto.OperatorID
	}
	if dto.DonePieces != nil {
		order.DonePieces = *dto.DonePieces
	}
	if dto.EstimatedMinutes != nil {
		order.EstimatedMinutes = *dto.EstimatedMinutes
	}
	if dto.Deadline != nil {
		deadlineTime, err := time.Parse("2006-01-02", *dto.Deadline)
		if err != nil {
			return fmt.Errorf("invalid deadline format: %w", err)
		}
		order.Deadline = deadlineTime
	}

	if dto.Items != nil {
		currentTotalDone := 0
		currentTotalPieces := 0
		totalPrice := 0.0

		// Liberar/asignar máquinas comparando con los items actuales
		for _, itemDTO := range *dto.Items {
			if itemDTO.MachineID != nil {
				oldMachineID := findOldMachineID(order.Items, itemDTO.ID)
				if oldMachineID != nil && *oldMachineID != *itemDTO.MachineID {
					db.Model(&types.Machine{}).Where("id = ?", *oldMachineID).Update("status", "idle")
				}
				if *itemDTO.MachineID != 0 {
					db.Model(&types.Machine{}).Where("id = ?", *itemDTO.MachineID).Update("status", "printing")
				}
			}
			currentTotalDone += itemDTO.DonePieces
			currentTotalPieces += itemDTO.Quantity
			totalPrice += itemDTO.Price
		}

		// OrderItem es una entidad débil: eliminar los existentes y reinsertar
		// Esto evita inconsistencias por IDs parciales o items huérfanos
		if err := db.Where("order_id = ?", order.Id).Delete(&types.OrderItem{}).Error; err != nil {
			return fmt.Errorf("failed to clear existing items: %w", err)
		}

		var newItems []types.OrderItem
		for _, itemDTO := range *dto.Items {
			price := itemDTO.Price
			newItems = append(newItems, types.OrderItem{
				OrderID:    order.Id,
				StlName:    itemDTO.StlName,
				Quantity:   itemDTO.Quantity,
				DonePieces: itemDTO.DonePieces,
				MaterialID: itemDTO.MaterialID,
				MachineID:  itemDTO.MachineID,
				Price:      &price,
			})
		}

		if len(newItems) > 0 {
			if err := db.Create(&newItems).Error; err != nil {
				return fmt.Errorf("failed to create new items: %w", err)
			}
		}

		order.DonePieces = currentTotalDone
		order.TotalPieces = currentTotalPieces
		order.TotalPrice = &totalPrice
	}

	// Calcular status automático si no fue completada manualmente
	if order.Status != "completed" {
		if order.DonePieces >= order.TotalPieces && order.TotalPieces > 0 {
			order.Status = "ready"
		} else if order.DonePieces > 0 {
			order.Status = "in-progress"
		}
	} else {
		// Si se marcó como completed, liberar todas las máquinas
		var currentItems []types.OrderItem
		db.Where("order_id = ?", order.Id).Find(&currentItems)
		for _, item := range currentItems {
			if item.MachineID != nil {
				newStatus := "idle"
				updmachine := types.UpdateMachineDTO{Status: &newStatus}
				servicesMachine.UpdateMachineUseCase(*item.MachineID, updmachine, companyID)
			}
		}
		now := time.Now()
		order.FinishTime = &now
	}

	// Guardar solo el order (los items ya están guardados)
	order.Items = nil
	if err := db.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	servicesWeb.InvalidateCache("orders:list:*")
	servicesWeb.PublishEvent("dashboard_updates", `{"type": "ORDER_UPDATED", "message": "ORDER UPDATED"}`)

	return nil
}

// findOldMachineID busca el MachineID previo de un item por su ID en la lista actual
func findOldMachineID(existingItems []types.OrderItem, itemID uint) *int {
	for _, existing := range existingItems {
		if existing.ID == itemID {
			return existing.MachineID
		}
	}
	return nil
}
