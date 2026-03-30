package services

import (
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetDashboardDataUseCase(userRole string, companyID uint) (*types.ProductionDashboardResponse, error) {
	db := storage.DatabaseInstance{}.Instance()
	var response types.ProductionDashboardResponse

	// --- 1. MACHINES ---
	var machines []types.Machine
	if err := db.Where("id_company = ?", companyID).Find(&machines).Error; err != nil {
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

	// --- 2. ACTIVE ORDERS ---
	var activeOrders []types.ProductionOrder

	err := db.Preload("Items").
		Preload("Items.Material").
		Preload("Items.Machine").
		Where("id_company = ? AND status IN ?", companyID, []string{"in-progress", "queued", "ready", "pending"}).
		Order("priority ASC").
		Find(&activeOrders).Error

	if err != nil {
		return &response, err
	}

	isAdmin := (userRole == "admin")

	if isAdmin {
		// --- 3. TOTAL MONTHLY REVENUE CALCULATION (Main Metric) ---
		var monthlyRevenue float64

		// Calculate the first day of the current month
		now := time.Now()
		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		err := db.Model(&types.ProductionOrder{}).
			Where("id_company = ? AND created_at >= ?", companyID, firstDayOfMonth).
			Select("COALESCE(SUM(price), 0)").
			Scan(&monthlyRevenue).Error

		if err != nil {
			monthlyRevenue = 0
		}

		response.TotalRevenueFDM = monthlyRevenue
		response.Orders = activeOrders

	} else {
		response.TotalRevenueFDM = 0
		censoredOrders := make([]types.ProductionOrder, len(activeOrders))
		copy(censoredOrders, activeOrders)

		for i := range censoredOrders {
			censoredOrders[i].TotalPrice = nil
		}
		response.Orders = censoredOrders
	}

	response.ActiveJobs = int64(len(activeOrders))

	return &response, nil
}
