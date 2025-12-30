package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func CreateMaterialUseCase(new types.CreateMaterialDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	if new.Name == "" || new.Type == "" {
		return fmt.Errorf("Name and Type are required")
	}

	var existmaterial types.Material
	err := db.Where("LOWER(name) = ?", strings.ToLower(new.Name)).First(&existmaterial).Error
	if err == nil {
		return fmt.Errorf("The Material Already Exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	newMaterial := types.Material{
		Name:        new.Name,
		Type:        new.Type,
		Description: new.Description,
		Brand:       new.Brand,
	}
	if err := db.Create(&newMaterial).Error; err != nil {
		return fmt.Errorf("The Product already exists")
	}
	services.InvalidateCache("materials:list:*")
	return nil
}
