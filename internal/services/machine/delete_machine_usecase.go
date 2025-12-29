package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/database/machine_utils"
)

func DeleteMachineUseCase(id int) error {
	db := storage.DatabaseInstance{}.Instance()
	machine, err := machineutils.GetMachinebyID(id, db)
	if err != nil {
		return err
	}

	return db.Delete(&machine).Error
}
