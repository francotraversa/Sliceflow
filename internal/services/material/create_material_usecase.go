package services

import (
	"fmt"

	materialutils "github.com/francotraversa/Sliceflow/internal/infra/database/material_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateMaterialUseCase(new types.CreateMaterialDTO) error {
	if new.Name == "" || new.Type == "" {
		return fmt.Errorf("Name and Type are required")
	}

	machine, err := materialutils.GetMaterial(new)
	if err != nil {
		return err
	}

	if machine != nil {
		return fmt.Errorf("The Material %s already exists", new.Name)
	}
	newMaterial := types.Material{
		Name:        new.Name,
		Type:        new.Type,
		Description: new.Description,
		Brand:       new.Brand,
	}

	if err := db_utils.Create(&newMaterial); err != nil {
		return fmt.Errorf("Error Creating Machine")
	}
	services.InvalidateCache("materials:list:*")
	return nil
}
