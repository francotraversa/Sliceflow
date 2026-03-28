package services

import "github.com/francotraversa/Sliceflow/internal/types"

// OrderRepository defines the contract for order persistence operations.
// Any implementation (Postgres, SQLite, mock) can satisfy this interface.
type OrderRepository interface {
	FindByID(id uint) (*types.ProductionOrder, error)
	Create(order *types.ProductionOrder) error
}

// MachineService defines operations on machines that the order use case needs.
type MachineService interface {
	UpdateStatus(machineID int, status string, companyID uint) error
}

// EventBus defines the contract for cache invalidation and real-time events.
type EventBus interface {
	InvalidateCache(pattern string)
	PublishEvent(channel string, message string)
}
