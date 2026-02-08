package services

import (
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetStockUseCase(query string) (*[]types.StockItem, error) {
	var allProducts []types.StockItem
	cacheKey := "stock:list:all"

	// 1. Traer todo de cache (esto funciona, por eso ves 1ms)
	if !services.GetCache(cacheKey, &allProducts) {
		db := storage.DatabaseInstance{}.Instance()
		if err := db.Find(&allProducts).Error; err != nil {
			return nil, err
		}
		services.SetCache(cacheKey, &allProducts)
	}

	// 2. Si NO hay query, devolvemos todo
	if query == "" {
		return &allProducts, nil
	}

	// 3. FILTRADO (Acá es donde estaba fallando)
	var filtered []types.StockItem
	q := strings.ToLower(strings.TrimSpace(query)) // Limpiamos espacios y pasamos a minúsculas

	for _, p := range allProducts {
		// Normalizamos los datos del producto para comparar
		productName := strings.ToLower(p.Name)
		productSKU := strings.ToLower(p.SKU)

		// Verificamos si coincide el SKU exacto O si el nombre contiene la búsqueda
		if productSKU == q || strings.Contains(productName, q) {
			filtered = append(filtered, p)
		}
	}

	return &filtered, nil
}
