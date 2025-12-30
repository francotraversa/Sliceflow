package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/database/machine_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
)

func DeleteMachineUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	machine, err := machineutils.GetMachinebyID(id, db)
	if err != nil {
		return err
	}
	services.InvalidateCache("machine:list:*")
	return db.Delete(&machine).Error
}
