package controller

import (
	"fmt"
	"net/http"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/stock"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateProductHandler godoc
// @Summary      Create a new product
// @Description  Registers a product in the catalog using the SKU (barcode) and name.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        product  body      types.ProductCreateRequest  true  "Product data"
// @Success      200      {string}  string                      "Product has been created"
// @Failure      400      {string}  string                      "Invalid JSON or validation error"
// @Failure      409      {string}  string                      "SKU already exists"
// @Security BearerAuth
// @Router       /hornero/authed/stock/product [post]
func CreateProductHandler(c echo.Context) error {
	var item types.ProductCreateRequest

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	err = services.CreateProductUseCase(item, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The product %s has been created", item.Name)})

}

func GetProductsHandler(c echo.Context) error {
	search := c.QueryParam("q")
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	result, err := services.GetStockUseCase(search, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &result)
}

// DeleteIdProductHandler godoc
// @Summary      Soft-delete product
// @Description  Performs a soft-delete of the product using its SKU.
// @Tags         Stock
// @Param        sku      path      int                       true  "Product SKU"
// @Success      200      {string}  string                      "The Product has been deleted"
// @Failure      400      {string}  string                      "Invalid ID"
// @Failure      409      {string}  string                      "Delete error"
// @Security BearerAuth
// @Router       /hornero/authed/stock/{sku} [delete]
func DeleteIdProductHandler(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "SKU is required"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	err = services.DeleteByIdUseCase(sku, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The product %s has been deleted", sku)})
}

// UpdateProductHandler godoc
// @Summary      Update product data
// @Description  Edits name, description, status or minimum stock of an existing product.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        sku      path      int                         true  "Product SKU"
// @Param        product  body      types.ProductUpdateRequest  true  "Data to update"
// @Success      200      {string}  string                      "Product updated successfully"
// @Failure      400      {string}  string                      "Validation error or invalid ID"
// @Failure      404      {string}  string                      "Product not found"
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
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	product, err := services.UpdateByIdProductUseCase(sku, item, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The Product %s has been updated", product.Name)})
}

// AddMovementHandler godoc
// @Summary      Register stock entry or exit
// @Description  Creates a movement and updates stock automatically.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        movement body types.CreateMovementRequest true "Movement data"
// @Success      201      {string} string "Movement created successfully"
// @Failure      400      {string} string "Validation error or insufficient stock"
// @Router       /hornero/authed/stock/movement [post]
func CreateMovementHandler(c echo.Context) error {
	var mov types.CreateMovementRequest

	if err := c.Bind(&mov); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	mov.UserID = claims.UserId

	err = services.AddStockMovementUseCase(mov, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: ("The Movement has been registered successfully")})

}

// GetStockHistoryHandler godoc
// @Summary      Get filtered history
// @Description  Filter movements by SKU, start date and end date.
// @Tags         Stock
// @Accept       json
// @Produce      json
// @Param        sku         query     string  false  "Product SKU"
// @Param        start_date  query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End date (YYYY-MM-DD)"
// @Success      200         {array}   types.StockMovement
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /hornero/authed/stock/history [get]
func GetStockHistoryHandler(c echo.Context) error {
	var filter types.HistoryFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid query parameters"})

	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	history, err := services.GetStockHistoryUseCase(filter, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, history)
}

// GetDashboardHandler godoc
// @Summary      Get business metrics
// @Tags         Dashboard
// @Security     BearerAuth
// @Success      200  {object}  types.DashboardResponse
// @Router       /hornero/authed/stock/movement/dashboard [get]
func GetDashboardHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	stats, err := services.GetDashboardStatsUseCase(claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &stats)
}
