package materialutils

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func GetMaterialbyID(id int, companyID uint) (*types.Material, error) {
	db := storage.DatabaseInstance{}.Instance()
	var material types.Material

	if err := db.Where("id = ? AND id_company = ?", id, companyID).First(&material).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Material doesn't exists")
		}
		return nil, err
	}

	return &material, nil
}

func GetMaterial(dto types.CreateMaterialDTO) (*types.Material, error) {
	db := storage.DatabaseInstance{}.Instance()
	var material types.Material

	// Use First to find by name
	if err := db.Where("name = ?", dto.Name).First(&material).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error database lookup for machine %s: %w", dto.Name, err)
	}

	// 3. If we get here, the material exists.
	return &material, nil
}
