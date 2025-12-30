package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func CreateMachineUseCase(dto types.CreateMachineDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	var existing types.Machine
	err := db.Where("LOWER(name) = ?", strings.ToLower(dto.Name)).First(&existing).Error
	if err == nil {
		return fmt.Errorf("A machine with that name already exists.")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	newMachine := types.Machine{
		Name:   dto.Name,
		Type:   dto.Type,
		Status: "idle",
	}
	if err := db.Create(&newMachine).Error; err != nil {
		return fmt.Errorf("The Machine already exists")
	}
	services.InvalidateCache("machine:list:*")
	return nil
}
