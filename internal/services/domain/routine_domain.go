package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type RoutineRepository interface {
	GetActiveOrders() ([]types.ProductionOrder, error)
	BulkUpdateOrders(orders []types.ProductionOrder) error
}

type RoutineUseCase interface {
	CheckAndSetPriorities() error
}
