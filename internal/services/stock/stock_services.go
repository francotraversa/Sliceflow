package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
)

type StockService struct {
	repo domain.StockRepository
}

func NewStockService(repo domain.StockRepository) domain.StockService {
	return &StockService{repo: repo}
}
