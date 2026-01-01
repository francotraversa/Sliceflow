package machineutils

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func GetMachinebyID(id int) (*types.Machine, error) {
	db := storage.DatabaseInstance{}.Instance()

	var machine types.Machine

	if err := db.First(&machine, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Machine doesn't exists")
		}
		return nil, err
	}

	return &machine, nil
}

func GetMachine(dto types.CreateMachineDTO) (*types.Machine, error) {
	db := storage.DatabaseInstance{}.Instance()
	var machine types.Machine

	// Usamos First
	if err := db.Where("name = ?", dto.Name).First(&machine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error database lookup for machine %s: %w", dto.Name, err)
	}

	// 3. Si llegamos acá, la máquina existe.
	return &machine, nil
}
