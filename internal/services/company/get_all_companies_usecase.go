package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllCompaniesUseCase() ([]types.Company, error) {
	db := storage.DatabaseInstance{}.Instance()

	var companies []types.Company
	if err := db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}
