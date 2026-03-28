package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func DeleteCompanyUseCase(IdCompany string) error {
	db := storage.DatabaseInstance{}.Instance()
	if err := db.Where("id_company = ?", IdCompany).Delete(&types.Company{}).Error; err != nil {
		return err
	}
	return nil
}
