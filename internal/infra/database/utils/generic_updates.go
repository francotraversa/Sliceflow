package db_utils

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
)

func Save[T any](entity *T) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Save(entity).Error; err != nil {
		return err
	}

	return nil
}
func Create[T any](entity *T) error {
	db := storage.DatabaseInstance{}.Instance()

	if err := db.Create(entity).Error; err != nil {
		return fmt.Errorf("error creating %T: %w", entity, err)
	}

	return nil
}
