package repository

import (
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type routineRepository struct {
	db *gorm.DB
}

func NewRoutineRepository(db *gorm.DB) *routineRepository {
	return &routineRepository{db: db}
}

func (r *routineRepository) GetActiveOrders() ([]types.ProductionOrder, error) {
	var orders []types.ProductionOrder
	if err := r.db.Where("status IN ?", []string{"pending", "queued", "in-progress"}).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *routineRepository) BulkUpdateOrders(orders []types.ProductionOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range orders {
			if err := tx.Save(&order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
