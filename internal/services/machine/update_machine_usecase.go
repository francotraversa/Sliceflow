package services

import (
	"fmt"

	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateMachineUseCase(id int, dto types.UpdateMachineDTO) error {

	machine, err := machineutils.GetMachinebyID(id)
	if err != nil {
		return err
	}

	// 2. Actualizar campos
	machine.Name = *dto.Name
	machine.Type = *dto.Type
	machine.Status = *dto.Status // Importante para cambiar estado manual

	if err := db_utils.Save(machine); err != nil {
		return fmt.Errorf("The Machine was not updated")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
