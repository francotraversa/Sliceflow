package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetMetricsUseCase(companyID uint) (*types.MetricsResponse, error) {
	db := storage.DatabaseInstance{}.Instance()
	var response types.MetricsResponse

	// --- 1. Fetch all machines for the company ---
	var machines []types.Machine
	if err := db.Where("id_company = ?", companyID).Find(&machines).Error; err != nil {
		return &response, err
	}

	machineMap := make(map[uint]types.MachineMetric)
	for _, m := range machines {
		machineMap[uint(m.ID)] = types.MachineMetric{
			MachineID:   m.ID,
			MachineName: m.Name,
			QueuedHours: 0,
		}
	}

	// --- 2. Fetch all materials for the company ---
	var materials []types.Material
	if err := db.Where("id_company = ?", companyID).Find(&materials).Error; err != nil {
		return &response, err
	}

	materialMap := make(map[uint]types.MaterialMetric)
	for _, mat := range materials {
		materialMap[uint(mat.ID)] = types.MaterialMetric{
			MaterialID:   mat.ID,
			MaterialName: mat.Name,
			MaterialType: mat.Type,
			QueuedKilos:  0,
		}
	}

	// --- 3. Fetch active orders ---
	var activeOrders []types.ProductionOrder
	err := db.Preload("Items").
		Where("id_company = ? AND status IN ?", companyID, []string{"in-progress", "queued", "pending"}).
		Find(&activeOrders).Error

	if err != nil {
		return &response, err
	}

	// --- 4. Calculate metrics ---
	for _, order := range activeOrders {
		for _, item := range order.Items {
			remainingPieces := item.Quantity - item.DonePieces
			if remainingPieces <= 0 {
				continue
			}

			if item.MachineID != nil && item.Time != nil {
				if val, ok := machineMap[uint(*item.MachineID)]; ok {
					val.QueuedHours += float64(remainingPieces*(*item.Time)) / 60.0
					machineMap[uint(*item.MachineID)] = val
				}
			}

			if item.MaterialID != nil && item.Weight != nil {
				if val, ok := materialMap[uint(*item.MaterialID)]; ok {
					val.QueuedKilos += float64(remainingPieces) * (*item.Weight) / 1000.0
					materialMap[uint(*item.MaterialID)] = val
				}
			}
		}
	}

	// --- 5. Assemble response ---
	for _, m := range machineMap {
		response.Machines = append(response.Machines, m)
	}

	for _, m := range materialMap {
		response.Materials = append(response.Materials, m)
	}

	// Ensure empty arrays instead of null in JSON if no elements
	if response.Machines == nil {
		response.Machines = []types.MachineMetric{}
	}
	if response.Materials == nil {
		response.Materials = []types.MaterialMetric{}
	}

	return &response, nil
}
