package services

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/database"
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

	return nil
}
