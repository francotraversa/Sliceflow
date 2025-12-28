package services

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func DeleteByIdUseCase(sku uint) error {
	db := storage.DatabaseInstance{}.Instance()

	result := db.Delete(&types.StockItem{}, sku)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("producto no encontrado o ya eliminado")
	}

	return nil
}
