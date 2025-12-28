package stock_utils

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func CheckProductExistsBySKU(sku string) (*types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	var item types.StockItem

	err := db.Where("sku = ?", sku).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}
