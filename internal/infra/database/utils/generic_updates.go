package db_utils

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
)

func Save[T any](entity *T, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Where("id_company = ?", companyID).Save(entity).Error; err != nil {
		return err
	}

	return nil
}

func Create[T any](entity *T) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Create(entity).Error; err != nil {
		return err
	}

	return nil
}

func CreateWithoutCompany[T any](entity *T) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Create(entity).Error; err != nil {
		return err
	}

	return nil
}

func SaveWithoutCompany[T any](entity *T) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Save(entity).Error; err != nil {
		return err
	}

	return nil
}
