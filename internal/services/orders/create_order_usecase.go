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
	order, err := ordersutils.CheckOrder(dto)
	if err != nil {
		return err
	}

	if order != nil {
		return fmt.Errorf("The Order %d already exists", dto.ID)
	}
	// 1. Validar que el Material exista (Integridad referencial)
	_, err = materialutils.GetMaterialbyID(dto.MaterialID)
	if err != nil {
		return err
	}

	_, err = machineutils.GetMachinebyID(*dto.MachineID)
	if err != nil {
		return err
	}
	// 3. Convertir Fecha (String "2025-12-31" -> Time)
	deadlineTime, err := time.Parse("2006-01-02", dto.Deadline)
	if err != nil {
		return errors.New("Format Date invalid (use YYYY-MM-DD)")
	}

	totalMinutes := (dto.EstimatedHours * 60) + dto.EstimatedMinutes

	initialStatus := "pending"
	if dto.MachineID != nil {
		initialStatus = "queued" // Si ya tiene máquina, pasa a cola
	}

	// 6. Armar el Modelo DB
	var machineID *int
	if dto.MachineID != nil {
		id := int(*dto.MachineID)
		machineID = &id
	}

	newOrder := types.ProductionOrder{
		ClientName:     dto.ClientName,
		ProductDetails: dto.ProductDetails,
		TotalPieces:    dto.TotalPieces,
		DonePieces:     0, // Arranca en 0

		MaterialID: dto.MaterialID, // Relación

		Priority:         dto.Priority,
		Notes:            dto.Notes,
		EstimatedMinutes: totalMinutes,
		Deadline:         deadlineTime,

		Status: initialStatus,

		OperatorID: dto.OperatorID,
		MachineID:  machineID, // Puntero (puede ser nil)
		Price:      &dto.Price,
	}

	// 7. Guardar

	if err := db_utils.Create(&newOrder); err != nil {
		return fmt.Errorf("Error Creating Machine")
	}
	services.InvalidateCache("orders:list:*")
	services.PublishEvent("dashboard_updates", `{"type": "ORDER_CREATED", "message": "NEW ORDER CREATED"}`)

	return nil
}
