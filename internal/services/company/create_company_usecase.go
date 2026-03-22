package services

import (
	"fmt"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateCompanyUseCase(company types.CompanyCreateDTO) error {
	db := storage.DatabaseInstance{}.Instance()

	NewCompany := types.Company{
		Name:      company.Name,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&NewCompany).Error; err != nil {
		return fmt.Errorf("error creating company: %v", err)
	}
	return nil
}
