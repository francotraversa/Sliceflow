package machineutils

import (
	"errors"

	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func GetMachinebyID(id int, db *gorm.DB) (*types.Machine, error) {
	var machine types.Machine

	if err := db.First(&machine, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Machine doesn't exists")
		}
		return nil, err
	}

	return &machine, nil
}
