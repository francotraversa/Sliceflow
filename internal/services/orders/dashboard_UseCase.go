package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetDashboardDataUseCase(userRole string) (*types.ProductionDashboardResponse, error) {
	db := storage.DatabaseInstance{}.Instance()
	var response types.ProductionDashboardResponse

	// --- 1. MÁQUINAS ---
	// Todos ven el estado de las máquinas
	var machines []types.Machine
	if err := db.Find(&machines).Error; err != nil {
		return &response, err
	}
	response.Machines = machines

	// Calcular Tasa de Utilización
	var busyMachines float64
	for _, m := range machines {
		if m.Status != "idle" && m.Status != "maintenance" {
			busyMachines++
		}
	}
	if len(machines) > 0 {
		response.UtilizationRate = (busyMachines / float64(len(machines))) * 100
	}

	// --- 2. ÓRDENES ACTIVAS ---
	var activeOrders []types.ProductionOrder

	err := db.Preload("Material").
		Preload("Machine").
		Where("status IN ?", []string{"in-progress", "queued"}).
		Order("priority ASC").
		Find(&activeOrders).Error

	if err != nil {
		return &response, err
	}
	response.Orders = activeOrders
	response.ActiveJobs = int64(len(activeOrders))

	isAdmin := (userRole == "admin")

	if isAdmin {
		var totalRevenue float64
		for _, o := range activeOrders {
			totalRevenue += o.Price
		}
		response.TotalRevenueFDM = totalRevenue

		response.Orders = activeOrders

	} else {
		// --- USUARIO NORMAL: CENSURA TOTAL ---

		// 1. El Total es Cero
		response.TotalRevenueFDM = 0
		response.TotalRevenueSLS = 0

		censoredOrders := make([]types.ProductionOrder, len(activeOrders))
		copy(censoredOrders, activeOrders)

		for i := range censoredOrders {
			censoredOrders[i].Price = 0
		}
		response.Orders = censoredOrders
	}

	response.ActiveJobs = int64(len(activeOrders))

	return &response, nil
}
