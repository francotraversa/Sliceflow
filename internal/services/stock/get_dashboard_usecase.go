package services

import (
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetDashboardStatsUseCase() (*types.DashboardResponse, error) {
	db := storage.DatabaseInstance{}.Instance()
	var response types.DashboardResponse

	db.Model(&types.StockItem{}).Where("status = ?", "active").Count(&response.TotalItems)

	db.Model(&types.StockItem{}).
		Where("status = ?", "active").
		Select("COALESCE(SUM(quantity * price), 0)").
		Scan(&response.TotalValue)

	db.Where("quantity <= min_qty AND status = ?", "active").
		Find(&response.LowStockItems)

	response.LowStockCount = int64(len(response.LowStockItems))

	startOfDay := time.Now().Truncate(24 * time.Hour)

	db.Model(&types.StockMovement{}).
		Where("created_at >= ?", startOfDay).
		Count(&response.MovementsToday)

	db.Model(&types.User{}).Count(&response.ActiveUsers)

	type Result struct {
		StockSKU  string
		TotalSold int
	}
	var results []Result

	err := db.Table("stock_movements").
		Select("stock_sku, SUM(qty_delta) as total_sold").
		Where("type = ?", "OUT").
		Group("stock_sku").
		Order("total_sold ASC"). // Orden ascendente porque son números negativos (ej: -100)
		Scan(&results).Error

	if err != nil {
		return &response, err
	}

	for _, res := range results {
		var item types.StockItem
		db.Select("Name").First(&item, "sku = ?", res.StockSKU)

		response.TopSellingItems = append(response.TopSellingItems, types.TopProduct{
			SKU:       res.StockSKU,
			Name:      item.Name,
			TotalSold: res.TotalSold * -1,
		})
	}

	return &response, nil
}
