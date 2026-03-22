package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	materialutils "github.com/francotraversa/Sliceflow/internal/infra/database/material_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
)

func DeleteMaterialUseCase(id int, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()
	material, err := materialutils.GetMaterialbyID(id, companyID)
	if err != nil {
		return err
	}
	services.InvalidateCache("materials:list:*")
	return db.Where("id = ? AND id_company = ?", id, companyID).Delete(&material).Error
}
