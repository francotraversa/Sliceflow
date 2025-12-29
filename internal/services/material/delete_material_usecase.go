package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	materialutils "github.com/francotraversa/Sliceflow/internal/database/material_utils"
)

func DeleteMaterialUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	material, err := materialutils.GetMaterialbyID(id, db)
	if err != nil {
		return err
	}
	return db.Delete(&material).Error
}
