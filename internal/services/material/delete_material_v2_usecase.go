package services

func (s *materialService) DeleteMaterial(id uint, companyID uint) error {
	return s.repo.Delete(id, companyID)
}
