package services

import (
	"fmt"

	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateMachineUseCase(dto types.CreateMachineDTO) error {
	machine, err := machineutils.GetMachine(dto)
	if err != nil {
		return err
	}

	if machine != nil {
		return fmt.Errorf("The machine %s already exists", dto.Name)
	}

	newMachine := types.Machine{
		Name:   dto.Name,
		Type:   dto.Type,
		Status: "idle",
	}
	if err := db_utils.Create(&newMachine); err != nil {
		return fmt.Errorf("Error Creating Machine")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
