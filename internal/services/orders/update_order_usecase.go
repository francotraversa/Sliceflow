package services

import (
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateOrderUseCase(id int, dto types.UpdateOrderDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	var order types.ProductionOrder

	// 1. Buscar la orden existente
	if err := db.First(&order, id).Error; err != nil {
		return err
	}

	// 2. Actualizar Datos Básicos
	if dto.ClientName != nil {
		order.ClientName = *dto.ClientName
	}

	if dto.ProductDetails != nil {
		order.ProductDetails = *dto.ProductDetails
	}

	if dto.TotalPieces != nil {
		order.TotalPieces = *dto.TotalPieces
	}

	if dto.DonePieces != nil {
		order.DonePieces = *dto.DonePieces
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

	// --- RELACIONES ---

	if dto.OperatorID != nil {
		order.OperatorID = *dto.OperatorID
	}

	if dto.MaterialID != nil {
		order.MaterialID = *dto.MaterialID
	}

	if dto.MachineID != nil {
		order.MachineID = dto.MachineID
	}

	// 4. Recalcular Tiempo (si mandaron datos)
	if dto.EstimatedHours != nil || dto.EstimatedMinutes != nil {
		hours := 0
		minutes := 0
		if dto.EstimatedHours != nil {
			hours = *dto.EstimatedHours
		}

		if dto.EstimatedMinutes != nil {
			minutes = *dto.EstimatedMinutes
		}

		order.EstimatedMinutes = (hours * 60) + minutes
	}

	if dto.Deadline != nil && *dto.Deadline != "" {
		parsedDeadline, err := time.Parse("2006-01-02", *dto.Deadline)
		if err == nil {
			order.Deadline = parsedDeadline
		}
	}

	if order.DonePieces >= order.TotalPieces && order.TotalPieces > 0 {
		order.Status = "completed"
	}

	if err := db_utils.Save(&order); err != nil {
		return fmt.Errorf("The Order was not updated")
	}
	services.InvalidateCache("orders:list:*")
	services.PublishEvent("dashboard_updates", `{"type": "ORDER_UPDATED", "message": "ORDER UPDATED"}`)

	return nil
}
