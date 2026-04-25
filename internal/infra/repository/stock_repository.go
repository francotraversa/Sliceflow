package repository

import (
	"errors"
	"time"

	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

// Delete implements [domain.StockRepository].
func (r *stockRepository) Delete(id uint, companyID uint) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Delete(&types.StockItem{}).Error
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(stock *types.StockItem) error {
	return r.db.Create(stock).Error
}

func (r *stockRepository) GetByID(id *uint, sku *string, companyID uint) (*types.StockItem, error) {
	var stock types.StockItem
	var err error
	if id != nil && sku == nil {
		err = r.db.Where("id = ? AND id_company = ?", *id, companyID).First(&stock).Error
	} else if id == nil && sku != nil {
		err = r.db.Where("sku = ? AND id_company = ?", *sku, companyID).First(&stock).Error
	} else if id != nil && sku != nil {
		err = r.db.Where("id = ? AND sku = ? AND id_company = ?", *id, *sku, companyID).First(&stock).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) GetAll(companyID uint) (*[]types.StockItem, error) {
	var stocks []types.StockItem
	err := r.db.Where("id_company = ?", companyID).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return &stocks, nil
}

func (r *stockRepository) Update(id uint, stock *types.StockItem, companyID uint) error {
	return r.db.Where("id = ? AND sku = ? AND id_company = ?", id, stock.SKU, companyID).Updates(stock).Error
}

func (r *stockRepository) GetDashboardStats(companyID uint) (*types.DashboardResponse, error) {
	var response types.DashboardResponse

	// Total active items count
	if err := r.db.Model(&types.StockItem{}).
		Where("status = ? AND id_company = ?", "active", companyID).
		Count(&response.TotalItems).Error; err != nil {
		return nil, err
	}

	// Total inventory value
	if err := r.db.Model(&types.StockItem{}).
		Where("status = ? AND id_company = ?", "active", companyID).
		Select("COALESCE(SUM(quantity * price), 0)").
		Scan(&response.TotalValue).Error; err != nil {
		return nil, err
	}

	// Low stock items (scoped to company)
	if err := r.db.Where("quantity <= min_qty AND status = ? AND id_company = ?", "active", companyID).
		Find(&response.LowStockItems).Error; err != nil {
		return nil, err
	}
	response.LowStockCount = int64(len(response.LowStockItems))

	// Movements today (scoped to company)
	startOfDay := time.Now().Truncate(24 * time.Hour)
	if err := r.db.Model(&types.StockMovement{}).
		Where("created_at >= ? AND id_company = ?", startOfDay, companyID).
		Count(&response.MovementsToday).Error; err != nil {
		return nil, err
	}

	// Active users in company
	if err := r.db.Model(&types.User{}).
		Where("id_company = ?", companyID).
		Count(&response.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Top selling items (scoped to company via stock_movements.id_company)
	type result struct {
		StockSKU  string
		TotalSold int
	}
	var results []result
	if err := r.db.Table("stock_movements").
		Select("stock_sku, SUM(qty_delta) as total_sold").
		Where("type = ? AND id_company = ?", "OUT", companyID).
		Group("stock_sku").
		Order("total_sold ASC"). // ASC porque los deltas OUT son negativos
		Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, res := range results {
		var item types.StockItem
		r.db.Select("name").First(&item, "sku = ? AND id_company = ?", res.StockSKU, companyID)
		response.TopSellingItems = append(response.TopSellingItems, types.TopProduct{
			SKU:       res.StockSKU,
			Name:      item.Name,
			TotalSold: res.TotalSold * -1,
		})
	}

	return &response, nil
}
