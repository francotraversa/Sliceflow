package controller

import (
	"fmt"
	"net/http"
	"strconv"

	services "github.com/francotraversa/Sliceflow/internal/services/stock"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateProductHandler godoc
// @Summary      Crear un nuevo producto
// @Description  Registra un producto en el catálogo usando el SKU (código de barras) y nombre.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        product  body      types.ProductCreateRequest  true  "Datos del producto a crear"
// @Success      200      {string}  string                      "Product has been created"
// @Failure      400      {string}  string                      "Invalid Json o error de validación"
// @Failure      409      {string}  string                      "El SKU ya existe"
// @Security     ApiKeyAuth
// @Router       /hornero/loged/stock/product [post]
func CreateProductHandler(c echo.Context) error {
	var item types.ProductCreateRequest

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}

	err := services.CreateProductUseCase(item)
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	return c.JSON(http.StatusOK, "Product has been created")

}

// GetAllProductsHandler godoc
// @Summary      Listar todos los productos
// @Description  Obtiene la lista completa de productos en stock que no han sido borrados.
// @Tags         Stock
// @Produce      json
// @Success      200      {array}   types.StockItem
// @Failure      409      {string}  string                      "Error al obtener productos"
// @Security     ApiKeyAuth
// @Router       /hornero/loged/stock/list [get]
func GetAllProductsHandler(c echo.Context) error {
	items, err := services.GetAllProductsUseCase()
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	return c.JSON(http.StatusOK, items)

}

// GetIdProductHandler godoc
// @Summary      Obtener producto por ID
// @Description  Busca un producto específico usando su ID numérico de base de datos.
// @Tags         Stock
// @Produce      json
// @Param        sku      path      int                       true  "ID del producto"
// @Success      200      {object}  types.StockItem
// @Failure      400      {string}  string                      "ID inválido"
// @Failure      409      {string}  string                      "Producto no encontrado"
// @Security     ApiKeyAuth
// @Router       /hornero/loged/stock/{sku} [get]
func GetIdProductHandler(c echo.Context) error {
	idParam := c.Param("sku")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "SKU must be a valid number")
	}
	sku := uint(id64)
	item, err := services.GetByIdUseCase(sku)
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	return c.JSON(http.StatusOK, item)
}

// DeleteIdProductHandler godoc
// @Summary      Baja lógica de producto
// @Description  Realiza un soft-delete del producto usando su ID.
// @Tags         Stock
// @Param        sku      path      int                       true  "ID del producto"
// @Success      200      {string}  string                      "The Product has been deleted"
// @Failure      400      {string}  string                      "ID inválido"
// @Failure      409      {string}  string                      "Error al eliminar"
// @Security     ApiKeyAuth
// @Router       /hornero/loged/stock/{sku} [delete]
func DeleteIdProductHandler(c echo.Context) error {
	idParam := c.Param("sku")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "SKU must be a valid number")
	}
	sku := uint(id64)
	err = services.DeleteByIdUseCase(sku)
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	return c.JSON(http.StatusOK, "The Product has been delated")
}

// UpdateProductHandler godoc
// @Summary      Actualizar datos de un producto
// @Description  Edita nombre, descripción, status o stock mínimo de un producto existente.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        sku      path      int                         true  "ID del producto"
// @Param        product  body      types.ProductUpdateRequest  true  "Datos a actualizar"
// @Success      200      {string}  string                      "Product updated successfully"
// @Failure      400      {string}  string                      "Error de validación o ID inválido"
// @Failure      404      {string}  string                      "Producto no encontrado"
// @Security     ApiKeyAuth
// @Router       /hornero/loged/stock/{sku} [put]
func UpdateByIdProductHandler(c echo.Context) error {
	sku := c.Param("sku")

	var item types.ProductUpdateRequest

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}

	product, err := services.UpdateByIdProductUseCase(sku, item)
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("The Product %s has been updated", product.SKU))
}
