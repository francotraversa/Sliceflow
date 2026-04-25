package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type MetricsRepository interface {
	GetMetrics(companyID uint) (*types.MetricsResponse, error)
}

type MetricsUseCase interface {
	GetMetrics(companyID uint) (*types.MetricsResponse, error)
}
