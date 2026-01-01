package services

import (
	"fmt"

	materialutils "github.com/francotraversa/Sliceflow/internal/infra/database/material_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateMaterialUseCase(id int, mat types.UpdateMaterialDTO) error {
	material, err := materialutils.GetMaterialbyID(id)
	if err != nil {
		return err
	}
	material.Name = mat.Name
	material.Type = mat.Type
	material.Description = mat.Description
	material.Brand = mat.Brand
	if err := db_utils.Save(material); err != nil {
		return fmt.Errorf("Error update Material %s", mat.Name)
	}
	services.InvalidateCache("materials:list:*")
	return nil
}
