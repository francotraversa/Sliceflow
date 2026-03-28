package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetStockUseCase(query string, companyID uint) (*[]types.StockItem, error) {
	var allProducts []types.StockItem
	cacheKey := fmt.Sprintf("stock:list:%d", companyID)

	if !services.GetCache(cacheKey, &allProducts) {
		db := storage.DatabaseInstance{}.Instance()
		if err := db.Where("id_company = ?", companyID).Find(&allProducts).Error; err != nil {
			return nil, err
		}
		services.SetCache(cacheKey, &allProducts)
	}

	// 2. If there's no query, return everything
	if query == "" {
		return &allProducts, nil
	}

	// 3. FILTERING
	var filtered []types.StockItem
	q := strings.ToLower(strings.TrimSpace(query)) // Trim spaces and convert to lowercase

	for _, p := range allProducts {
		// Normalize product data for comparison
		productName := strings.ToLower(p.Name)
		productSKU := strings.ToLower(p.SKU)

		// Check for exact SKU match OR if the name contains the search term
		if productSKU == q || strings.Contains(productName, q) {
			filtered = append(filtered, p)
		}
	}

	return &filtered, nil
}
