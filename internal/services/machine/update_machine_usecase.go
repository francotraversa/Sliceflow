package services

import (
	"fmt"

	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateMachineUseCase(id int, dto types.UpdateMachineDTO, companyID uint) error {

	machine, err := machineutils.GetMachinebyID(id, companyID)
	if err != nil {
		return err
	}

	if dto.Name != nil {
		machine.Name = *dto.Name
	}

	if dto.Type != nil {
		machine.Type = *dto.Type
	}

	if dto.Status != nil {
		machine.Status = *dto.Status
	}

	if err := db_utils.Save(machine, companyID); err != nil {
		return fmt.Errorf("The Machine was not updated")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
