package services

import (
	"fmt"
	"time"

	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *OrderService) UpdateOrder(id uint, dto types.UpdateOrderDTO, idCompany uint) error {
	order, err := s.repo.GetOrderWithItems(id, idCompany)
	if err != nil {
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

	var newItems []types.OrderItem

	if dto.Items != nil {
		currentTotalDone := 0
		currentTotalPieces := 0
		totalPrice := 0.0

		for _, itemDTO := range *dto.Items {
			if itemDTO.MachineID != nil {
				oldMachineID := findOldMachineID(order.Items, itemDTO.ID)
				if oldMachineID != nil && *oldMachineID != *itemDTO.MachineID {
					statusIdle := "idle"
					s.machineService.UpdateMachine(uint(*oldMachineID), types.UpdateMachineDTO{Status: &statusIdle}, idCompany)
				}
				if *itemDTO.MachineID != 0 {
					statusPrinting := "printing"
					s.machineService.UpdateMachine(uint(*itemDTO.MachineID), types.UpdateMachineDTO{Status: &statusPrinting}, idCompany)
				}
			}
			currentTotalDone += itemDTO.DonePieces
			currentTotalPieces += itemDTO.Quantity
			totalPrice += itemDTO.Price
		}

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
				Weight:     itemDTO.Weight,
				Time:       itemDTO.Time,
			})
		}
		order.DonePieces = currentTotalDone
		order.TotalPieces = currentTotalPieces
		order.TotalPrice = &totalPrice
	}

	if order.Status != "completed" {
		if order.DonePieces >= order.TotalPieces && order.TotalPieces > 0 {
			order.Status = "ready"
		} else if order.DonePieces > 0 {
			order.Status = "in-progress"
		}
	} else {
		for _, item := range order.Items {
			if item.MachineID != nil {
				newStatus := "idle"
				s.machineService.UpdateMachine(uint(*item.MachineID), types.UpdateMachineDTO{Status: &newStatus}, idCompany)
			}
		}
		now := time.Now()
		order.FinishTime = &now
	}

	if err := s.repo.UpdateFullOrder(order, newItems); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	servicesWeb.InvalidateCache("orders:list:*")
	servicesWeb.PublishEvent("dashboard_updates", `{"type": "ORDER_UPDATED", "message": "ORDER UPDATED"}`)

	return nil
}
