package ordersutils

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func CheckOrder(order types.CreateOrderDTO) (*types.ProductionOrder, error) {
	db := storage.DatabaseInstance{}.Instance()
	var orderexists types.ProductionOrder

	err := db.Where("id = ?", order.ID).First(&orderexists).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderexists, nil
}
