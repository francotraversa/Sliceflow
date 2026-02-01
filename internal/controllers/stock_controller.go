package controller

import (
	"fmt"
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/auth"
	services "github.com/francotraversa/Sliceflow/internal/services/stock"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
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
// @Security BearerAuth
// @Router       /hornero/authed/stock/product [post]
func CreateProductHandler(c echo.Context) error {
	var item types.ProductCreateRequest

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	err := services.CreateProductUseCase(item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The product %s has been created", item.Name)})

}

// GetAllProductsHandler godoc
// @Summary      Listar todos los productos
// @Description  Obtiene la lista completa de productos en stock que no han sido borrados.
// @Tags         Stock
// @Produce      json
// @Success      200      {array}   types.StockItem
// @Failure      409      {string}  string                      "Error al obtener productos"
// @Security BearerAuth
// @Router       /hornero/authed/stock/list [get]
func GetAllProductsHandler(c echo.Context) error {
	items, err := services.GetAllProductsUseCase()
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
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
// @Security BearerAuth
// @Router       /hornero/authed/stock/{sku} [get]
func GetIdProductHandler(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "SKU is required"})

	}
	item, err := services.GetByIdUseCase(sku)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
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
// @Security BearerAuth
// @Router       /hornero/authed/stock/{sku} [delete]
func DeleteIdProductHandler(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "SKU is required"})
	}
	err := services.DeleteByIdUseCase(sku)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The product %s has been deleted", sku)})
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
// @Security BearerAuth
// @Router       /hornero/authed/stock/{sku} [put]
func UpdateByIdProductHandler(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "SKU is required"})
	}

	var item types.ProductUpdateRequest

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	product, err := services.UpdateByIdProductUseCase(sku, item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The Product %s has been updated", product.Name)})
}

// AddMovementHandler godoc
// @Summary      Registrar entrada o salida de stock
// @Description  Genera un movimiento y actualiza el stock automáticamente.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        movement body types.CreateMovementRequest true "Datos del movimiento"
// @Success      201      {string} string "Movement created successfully"
// @Failure      400      {string} string "Error de validación o Stock insuficiente"
// @Router       /hornero/authed/stock/movement [post]
func CreateMovementHandler(c echo.Context) error {
	var mov types.CreateMovementRequest

	if err := c.Bind(&mov); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.JwtCustomClaims)
	mov.UserID = claims.UserId

	err := services.AddStockMovementUseCase(mov)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: ("The Movement has been registered successfully")})

}

// GetStockHistoryHandler godoc
// @Summary      Obtener historial con filtros
// @Description  Permite filtrar movimientos por SKU, fecha de inicio y fecha de fin.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        sku         query     string  false  "SKU del producto"
// @Param        start_date  query     string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200         {array}   types.StockMovement
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /hornero/authed/stock/history [get]
func GetStockHistoryHandler(c echo.Context) error {
	var filter types.HistoryFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid query parameters"})

	}

	history, err := services.GetStockHistoryUseCase(filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// GetDashboardHandler godoc
// @Summary      Obtener métricas del negocio
// @Tags         Dashboard
// @Security     BearerAuth
// @Success      200  {object}  types.DashboardResponse
// @Router       /hornero/authed/stock/movement/dashboard [get]
func GetDashboardHandler(c echo.Context) error {
	stats, err := services.GetDashboardStatsUseCase()
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &stats)
}

// GetProductByNameHandler godoc
// @Summary      Obtener producto por nombre
// @Description  Busca un producto específico usando su nombre.
// @Tags         Stock
// @Produce      json
// @Param        name      path      string                       true  "Nombre del producto"
// @Success      200      {object}  types.StockItem
// @Failure      400      {string}  string                      "Nombre inválido"
// @Failure      409      {string}  string                      "Producto no encontrado"
// @Security BearerAuth
// @Router       /hornero/authed/stock/list/{name} [get]
func GetProductByNameHandler(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Name is required"})
	}
	item, err := services.GetProductByNameUseCase(name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}
