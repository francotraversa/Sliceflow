package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
)

func DeleteMachineUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	machine, err := machineutils.GetMachinebyID(id)
	if err != nil {
		return err
	}
	services.InvalidateCache("machine:list:*")
	return db.Delete(&machine).Error
}
