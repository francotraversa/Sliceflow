package materialutils

import (
	"errors"

	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func GetMaterialbyID(id int, db *gorm.DB) (*types.Material, error) {
	var material types.Material

	if err := db.First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Material doesn't exists")
		}
		return nil, err
	}

	return &material, nil
}
