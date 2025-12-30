package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/database/machine_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateMachineUseCase(id int, dto types.UpdateMachineDTO) error {
	db := storage.DatabaseInstance{}.Instance()

	machine, err := machineutils.GetMachinebyID(id, db)
	if err != nil {
		return err
	}

	// 2. Actualizar campos
	machine.Name = dto.Name
	machine.Type = dto.Type
	machine.Status = dto.Status // Importante para cambiar estado manual

	if err := db.Save(machine).Error; err != nil {
		return fmt.Errorf("The Machine was not updated")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
