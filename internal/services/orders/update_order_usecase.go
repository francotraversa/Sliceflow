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

func UpdateOrderUseCase(id int, dto types.UpdateOrderDTO, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder
	if err := db.Preload("Items").Where("id_company = ? AND id_order = ?", companyID, id).First(&order).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

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
		var updatedItems []types.OrderItem
		currentTotalDone := 0
		currentTotalPieces := 0
		totalPrice := 0.0

		for _, itemDTO := range *dto.Items {
			// If the item has a new machine, release the previous machine
			if itemDTO.MachineID != nil {
				oldMachineID := findOldMachineID(order.Items, itemDTO.ID)
				if oldMachineID != nil && *oldMachineID != *itemDTO.MachineID {
					// Release the previous machine
					db.Model(&types.Machine{}).Where("id = ?", *oldMachineID).Update("status", "idle")
				}
				// Assign the new machine
				if *itemDTO.MachineID != 0 {
					db.Model(&types.Machine{}).Where("id = ?", *itemDTO.MachineID).Update("status", "printing")
				}
			}

			item := types.OrderItem{
				ID:         itemDTO.ID,
				OrderID:    order.Id,
				StlName:    itemDTO.StlName,
				Quantity:   itemDTO.Quantity,
				DonePieces: itemDTO.DonePieces,
				MaterialID: itemDTO.MaterialID,
				MachineID:  itemDTO.MachineID,
				Price:      itemDTO.Price,
			}
			updatedItems = append(updatedItems, item)
			currentTotalDone += item.DonePieces
			currentTotalPieces += item.Quantity
			if item.Price != nil {
				totalPrice += *item.Price
			}
		}
		order.Items = updatedItems
		order.DonePieces = currentTotalDone
		order.TotalPieces = currentTotalPieces
		order.TotalPrice = &totalPrice
	}

	if order.Status != "completed" {
		if order.DonePieces >= order.TotalPieces {
			order.Status = "ready"
		} else if order.DonePieces > 0 {
			order.Status = "in-progress"
		}
	} else {
		for _, item := range order.Items {
			if item.MachineID != nil {
				newStatus := "idle"
				updmachine := types.UpdateMachineDTO{Status: &newStatus}
				servicesMachine.UpdateMachineUseCase(*item.MachineID, updmachine, companyID)
			}
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

// findOldMachineID finds the previous MachineID of an item by its ID in the existing items list
func findOldMachineID(existingItems []types.OrderItem, itemID uint) *int {
	for _, existing := range existingItems {
		if existing.ID == itemID {
			return existing.MachineID
		}
	}
	return nil
}
