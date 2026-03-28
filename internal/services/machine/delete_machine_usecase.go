package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
)

func DeleteMachineUseCase(id int, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()
	machine, err := machineutils.GetMachinebyID(id, companyID)
	if err != nil {
		return err
	}

	if err := db.Where("id_company = ?", companyID).Delete(&machine).Error; err != nil {
		return fmt.Errorf("The Machine was not deleted")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
