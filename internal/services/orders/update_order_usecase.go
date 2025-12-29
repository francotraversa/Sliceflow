package services

import (
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/database"
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
	order.ClientName = dto.ClientName
	order.ProductDetails = dto.ProductDetails
	order.TotalPieces = dto.TotalPieces
	order.DonePieces = dto.DonePieces // Actualizar progreso
	order.Priority = dto.Priority
	order.Notes = dto.Notes
	order.Status = dto.Status
	order.OperatorID = dto.OperatorID
	order.Price = dto.Price

	// 3. Actualizar Relaciones (Validar si existen es opcional pero recomendado)
	order.MaterialID = dto.MaterialID
	order.MachineID = dto.MachineID

	// 4. Recalcular Tiempo (si mandaron datos)
	totalMinutes := (dto.EstimatedHours * 60) + dto.EstimatedMinutes
	order.EstimatedMinutes = totalMinutes

	// 5. Parsear Fecha Límite
	if dto.Deadline != "" {
		parsedDeadline, err := time.Parse("2006-01-02", dto.Deadline)
		if err == nil {
			order.Deadline = parsedDeadline
		}
	}

	// 6. Lógica de Estado Automático (Opcional)
	// Ejemplo: Si completó todas las piezas, pasar a 'completed'
	if order.DonePieces >= order.TotalPieces && order.TotalPieces > 0 {
		order.Status = "completed"
	}

	// 7. Guardar cambios
	if err := db.Save(order).Error; err != nil {
		return fmt.Errorf("The Order was not updated")
	}
	return nil
}
