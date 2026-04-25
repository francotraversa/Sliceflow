package services

import "github.com/francotraversa/Sliceflow/internal/types"

func (s *StockService) GetDashboard(companyID uint) (*types.DashboardResponse, error) {
	return s.repo.GetDashboardStats(companyID)
}
