package services

import "github.com/francotraversa/Sliceflow/internal/services/domain"

type OrderService struct {
	repo           domain.OrderRepository
	machineService domain.MachineUseCase
}

func NewOrderService(repo domain.OrderRepository, machineService domain.MachineUseCase) domain.OrderUseCase {
	return &OrderService{repo: repo, machineService: machineService}
}
