package services

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func DeleteByIdUseCase(sku string) error {
	db := storage.DatabaseInstance{}.Instance()

	result := db.Where("sku = ?", sku).Delete(&types.StockItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Coudnt find product or has been deleted")
	}
	services.InvalidateCache("stock:list:all")
	services.PublishEvent("dashboard_updates", `{"type": "PRODUCT_DELETED", "message": "PRODUCT DELETED"}`)

	return nil
}
