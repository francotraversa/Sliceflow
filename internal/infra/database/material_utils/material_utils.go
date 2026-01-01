package materialutils

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func GetMaterialbyID(id int) (*types.Material, error) {
	db := storage.DatabaseInstance{}.Instance()
	var material types.Material

	if err := db.First(&material, id).Error; err != nil {
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

	// Usamos First
	if err := db.Where("name = ?", dto.Name).First(&material).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error database lookup for machine %s: %w", dto.Name, err)
	}

	// 3. Si llegamos acá, la máquina existe.
	return &material, nil
}
