package services

import (
	"errors"
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/database/machine_utils"
	materialutils "github.com/francotraversa/Sliceflow/internal/database/material_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateOrderUseCase(dto types.CreateOrderDTO) error {
	db := storage.DatabaseInstance{}.Instance()

	// 1. Validar que el Material exista (Integridad referencial)
	_, err := materialutils.GetMaterialbyID(dto.MaterialID, db)
	if err != nil {
		return err
	}

	_, err = machineutils.GetMachinebyID(*dto.MachineID, db)
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
		Price:      dto.Price,
	}

	// 7. Guardar
	if err := db.Create(&newOrder).Error; err != nil {
		return fmt.Errorf("The Order already exists")
	}
	return nil
}
