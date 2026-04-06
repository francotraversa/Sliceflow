package services

func (s *OrderService) DeleteOrder(id uint, companyID uint) error {
	return s.repo.DeleteOrder(id, companyID)
}
