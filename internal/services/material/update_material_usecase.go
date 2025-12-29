package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	materialutils "github.com/francotraversa/Sliceflow/internal/database/material_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateMaterialUseCase(id int, mat types.UpdateMaterialDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	material, err := materialutils.GetMaterialbyID(id, db)
	if err != nil {
		return err
	}
	material.Name = mat.Name
	material.Type = mat.Type
	material.Description = mat.Description
	material.Brand = mat.Brand

	if err := db.Save(material).Error; err != nil {
		return fmt.Errorf("The Product was not updated")
	}
	return nil
}
