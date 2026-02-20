package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllOrdersUseCase(filter types.OrderFilter) (*[]types.ProductionOrder, error) {
	db := storage.DatabaseInstance{}.Instance()

	// Actualizamos la cacheKey para que incluya las fechas,
	// si no, te mostraría datos viejos al filtrar.
	cacheKey := fmt.Sprintf("orders:list:st_%s:id_%d:from_%s:to_%s:sort_%v",
		filter.Status, filter.ID, filter.FromDate, filter.ToDate, filter.SortPriority)

	var orders []types.ProductionOrder

	if services.GetCache(cacheKey, &orders) {
		return &orders, nil
	}

	query := db.Preload("Material").Preload("Machine").Preload("Items")

	// 1. Filtro por ID (Prioridad absoluta)
	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	} else {
		// 2. Filtro por Estado
		if filter.Status != "" {
			// Si el front pide 'pending', incluimos 'ready' para que no desaparezcan
			if filter.Status == "pending" {
				query = query.Where("status IN ?", []string{"pending", "in-progress", "ready"})
			} else {
				query = query.Where("status = ?", filter.Status)
			}
		}

		// 3. FILTRO POR RANGO DE FECHAS (Nuevo)
		if filter.FromDate != "" && filter.ToDate != "" {
			query = query.Where("created_at BETWEEN ? AND ?",
				filter.FromDate+" 00:00:00",
				filter.ToDate+" 23:59:59")
		} else if filter.FromDate != "" {
			query = query.Where("created_at >= ?", filter.FromDate+" 00:00:00")
		} else if filter.ToDate != "" {
			query = query.Where("created_at <= ?", filter.ToDate+" 23:59:59")
		}
	}

	// 4. Ordenamiento
	if filter.SortPriority {
		query = query.Order("priority ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	services.SetCache(cacheKey, &orders)
	return &orders, nil
}
