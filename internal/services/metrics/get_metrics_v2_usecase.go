package services

import (
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *metricsService) GetMetrics(companyID uint) (*types.MetricsResponse, error) {
	return s.repo.GetMetrics(companyID)
}
