package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
)

type metricsService struct {
	repo domain.MetricsRepository
}

func NewMetricsService(repo domain.MetricsRepository) domain.MetricsUseCase {
	return &metricsService{repo: repo}
}
