package services

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func DeleteByIdUseCase(sku string, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()

	result := db.Where("sku = ? AND id_company = ?", sku, companyID).Delete(&types.StockItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Coudnt find product or has been deleted")
	}
	services.InvalidateCache(fmt.Sprintf("stock:list:%d", companyID))
	services.PublishEvent("dashboard_updates", `{"type": "PRODUCT_DELETED", "message": "PRODUCT DELETED"}`)

	return nil
}
