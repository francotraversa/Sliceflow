package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetDashboardDataUseCase(userRole string) (*types.ProductionDashboardResponse, error) {
	db := storage.DatabaseInstance{}.Instance()
	var response types.ProductionDashboardResponse

	// --- 1. MÁQUINAS ---
	var machines []types.Machine
	if err := db.Find(&machines).Error; err != nil {
		return &response, err
	}
	response.Machines = machines

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

	// CAMBIO CLAVE: Agregamos Preload("Items") para traer la lista de piezas
	err := db.Preload("Items").
		Preload("Material").
		Preload("Machine").
		Where("status IN ?", []string{"in-progress", "queued"}).
		Order("priority ASC").
		Find(&activeOrders).Error

	if err != nil {
		return &response, err
	}

	isAdmin := (userRole == "admin")

	if isAdmin {
		var totalRevenue float64
		for _, o := range activeOrders {
			// CAMBIO: Usamos TotalPrice que es el campo del nuevo modelo
			totalRevenue += *o.Price
		}
		response.TotalRevenueFDM = totalRevenue
		response.Orders = activeOrders

	} else {
		// --- USUARIO NORMAL: CENSURA ---
		response.TotalRevenueFDM = 0
		response.TotalRevenueSLS = 0

		// Creamos una copia para no afectar la data original si fuera necesario
		censoredOrders := make([]types.ProductionOrder, len(activeOrders))
		copy(censoredOrders, activeOrders)

		for i := range censoredOrders {
			// CAMBIO: Ponemos el precio en 0 para usuarios no admin
			censoredOrders[i].Price = nil
		}
		response.Orders = censoredOrders
	}

	response.ActiveJobs = int64(len(activeOrders))

	return &response, nil
}
