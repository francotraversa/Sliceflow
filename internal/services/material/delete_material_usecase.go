package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	materialutils "github.com/francotraversa/Sliceflow/internal/infra/database/material_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
)

func DeleteMaterialUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	material, err := materialutils.GetMaterialbyID(id)
	if err != nil {
		return err
	}
	services.InvalidateCache("materials:list:*")
	return db.Delete(&material).Error
}
