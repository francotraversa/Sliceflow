package services

import (
	"errors"
	"fmt"
	"time"

	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	materialutils "github.com/francotraversa/Sliceflow/internal/infra/database/material_utils"
	ordersutils "github.com/francotraversa/Sliceflow/internal/infra/database/orders_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateOrderUseCase(dto types.CreateOrderDTO) error {
	if dto.ID != nil {
		order, _ := ordersutils.CheckOrder(dto)
		if order != nil {
			return fmt.Errorf("The Order %d already exists", *dto.ID)
		}
	}

	// 1. Validar que el Material exista (Integridad referencial)
	_, err := materialutils.GetMaterialbyID(dto.MaterialID)
	if err != nil {
		return err
	}

	if dto.MachineID != nil {
		_, err = machineutils.GetMachinebyID(*dto.MachineID)
		if err != nil {
			return fmt.Errorf("Machine could not be found: %w", err)
		}
	}
	var itemsDB []types.OrderItem
	totalPiecesCalculated := 0

	for _, itemDTO := range dto.Items {
		totalPiecesCalculated += itemDTO.Quantity

		itemsDB = append(itemsDB, types.OrderItem{
			StlName:    itemDTO.ProductName,
			Quantity:   itemDTO.Quantity,
			DonePieces: 0,
		})
	}
	deadlineTime, err := time.Parse("2006-01-02", dto.Deadline)
	if err != nil {
		return errors.New("Format Date invalid (use YYYY-MM-DD)")
	}

	totalMinutes := (dto.EstimatedHours * 60) + dto.EstimatedMinutes

	initialStatus := "pending"
	if dto.MachineID != nil {
		initialStatus = "queued" // Si ya tiene máquina, pasa a cola
	}

	newOrder := types.ProductionOrder{
		ID:               *dto.ID,
		ClientName:       dto.ClientName,
		Items:            itemsDB,               // La lista de piezas que armamos en el loop
		TotalPieces:      totalPiecesCalculated, // Usamos la suma de las cantidades
		DonePieces:       0,
		MaterialID:       dto.MaterialID,
		Priority:         dto.Priority,
		Notes:            dto.Notes,
		EstimatedMinutes: totalMinutes,
		Deadline:         deadlineTime,
		Status:           initialStatus,
		OperatorID:       dto.OperatorID,
		MachineID:        dto.MachineID,
		Price:            dto.Price, // Asignamos el precio del DTO al TotalPrice del modelo
	}

	if err := db_utils.Create(&newOrder); err != nil {
		return fmt.Errorf("could not save order and items: %w", err)
	}
	services.InvalidateCache("orders:list:*")
	services.PublishEvent("dashboard_updates", `{"type": "ORDER_CREATED", "message": "NEW MULTI-ITEM ORDER CREATED"}`)

	return nil
}
