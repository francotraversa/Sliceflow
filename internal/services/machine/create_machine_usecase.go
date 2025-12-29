package services

import (
	"errors"
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func CreateMachineUseCase(dto types.CreateMachineDTO) error {
	db := storage.DatabaseInstance{}.Instance()
	var existing types.Machine
	err := db.Where("LOWER(name) = ?", strings.ToLower(dto.Name)).First(&existing).Error
	if err == nil {
		return errors.New("ya existe una máquina con ese nombre")
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
	return nil
}
