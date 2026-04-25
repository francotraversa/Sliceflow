package repository

import (
	"fmt"
	"time"

	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *types.ProductionOrder) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint, companyID uint) (bool, error) {
	var order types.ProductionOrder
	err := r.db.Where("id_order = ? AND id_company = ?", id, companyID).First(&order).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *orderRepository) GetOrderWithItems(id uint, companyID uint) (*types.ProductionOrder, error) {
	var order types.ProductionOrder
	err := r.db.Preload("Items").Where("id_order = ? AND id_company = ?", id, companyID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateFullOrder(order *types.ProductionOrder, newItems []types.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if newItems != nil {
			if err := tx.Where("order_id = ?", order.Id).Delete(&types.OrderItem{}).Error; err != nil {
				return err
			}
			if len(newItems) > 0 {
				if err := tx.Create(&newItems).Error; err != nil {
					return err
				}
			}
		}
		// Save order
		order.Items = nil
		return tx.Save(order).Error
	})
}

func (r *orderRepository) GetOrdersByFilter(filter types.OrderFilter, companyID uint) (*[]types.ProductionOrder, error) {
	var orders []types.ProductionOrder
	query := r.db.Preload("Items.Material").Preload("Items.Machine").Preload("Items").Where("id_company = ?", companyID)

	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	} else {
		if filter.Status != "" {
			if filter.Status == "pending" {
				query = query.Where("status IN ?", []string{"pending", "in-progress", "ready"})
			} else {
				query = query.Where("status = ?", filter.Status)
			}
		}

		if filter.FromDate != "" && filter.ToDate != "" {
			query = query.Where("created_at BETWEEN ? AND ?",
				filter.FromDate+" 00:00:00",
				filter.ToDate+" 23:59:59")
		} else if filter.FromDate != "" {
			query = query.Where("created_at >= ?", filter.FromDate+" 00:00:00")
		} else if filter.ToDate != "" {
			query = query.Where("created_at <= ?", filter.ToDate+" 23:59:59")
		}
	}

	if filter.SortPriority {
		query = query.Order("priority ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

func (r *orderRepository) DeleteOrder(id uint, companyID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order types.ProductionOrder

		if err := tx.Preload("Items").Where("id_company = ? AND id = ?", companyID, id).First(&order).Error; err != nil {
			return fmt.Errorf("order not found: %w", err)
		}

		for _, item := range order.Items {
			if item.MachineID != nil && *item.MachineID != 0 {
				if err := tx.Model(&types.Machine{}).Where("id = ?", *item.MachineID).Update("status", "idle").Error; err != nil {
					return fmt.Errorf("failed to set machine to idle: %w", err)
				}
				services.PublishEvent("dashboard_updates", `{"type": "MACHINE_STATUS_CHANGED", "message": "Machine set to idle due to order deletion"}`)
			}
		}

		if err := tx.Delete(&order).Error; err != nil {
			return fmt.Errorf("failed to delete order: %w", err)
		}
		return nil
	})
}

func (r *orderRepository) DashboardOrders(companyID uint) (*types.ProductionDashboardResponse, error) {
	var response types.ProductionDashboardResponse

	var machines []types.Machine
	if err := r.db.Where("id_company = ?", companyID).Find(&machines).Error; err != nil {
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

	var activeOrders []types.ProductionOrder

	err := r.db.Preload("Items").
		Preload("Items.Material").
		Preload("Items.Machine").
		Where("id_company = ? AND status IN ?", companyID, []string{"in-progress", "queued", "ready", "pending"}).
		Order("priority ASC").
		Find(&activeOrders).Error

	if err != nil {
		return &response, err
	}

	response.Orders = activeOrders
	response.ActiveJobs = int64(len(activeOrders))

	var monthlyRevenue float64
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	errRev := r.db.Model(&types.ProductionOrder{}).
		Where("id_company = ? AND created_at >= ?", companyID, firstDayOfMonth).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&monthlyRevenue).Error
	if errRev == nil {
		response.TotalRevenueFDM = monthlyRevenue
	} else {
		response.TotalRevenueFDM = 0
	}

	return &response, nil
}
