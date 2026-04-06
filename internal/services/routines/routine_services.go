package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
)

type RoutineService struct {
	repo domain.RoutineRepository
}

func NewRoutineService(repo domain.RoutineRepository) domain.RoutineUseCase {
	return &RoutineService{repo: repo}
}
