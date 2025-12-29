package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
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

	// 👁️ TODOS VEN TODO: Eliminamos el filtro 'Where operator_id = ?'
	// Queremos que el operario vea la cola completa de trabajo.
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

	// --- 3. REVENUE (DINERO) - BLOQUEO DE SEGURIDAD 🔒 ---
	if userRole == "admin" {
		// SI ES ADMIN: Calculamos la plata real.
		var totalRevenue float64

		// Sumamos el precio de las órdenes activas como ejemplo
		for _, o := range activeOrders {
			totalRevenue += o.Price
		}

		response.TotalRevenueFDM = totalRevenue

	} else {
		// SI NO ES ADMIN: Se muestra $0.00
		// Ven el trabajo, pero no la facturación.
		response.TotalRevenueFDM = 0
		response.TotalRevenueSLS = 0
	}

	return &response, nil
}
