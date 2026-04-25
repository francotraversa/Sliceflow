package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

type StockController struct {
	stockService domain.StockService
}

func NewStockController(stockService domain.StockService) *StockController {
	return &StockController{stockService: stockService}
}

func (c *StockController) CreateItem(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	var stock types.ProductCreateRequest
	if err := ctx.Bind(&stock); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.stockService.CreateItem(&stock, claims.CompanyId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The product %s has been created", stock.Name)})
}

func (c *StockController) UpdateItem(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid stock ID")
	}

	var stock types.ProductUpdateRequest
	if err := ctx.Bind(&stock); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.stockService.UpdateItem(uint(id), &stock, claims.CompanyId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The product %s has been updated", stock.Name)})
}

func (c *StockController) DeleteItem(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid stock ID")
	}

	companyID := ctx.Get("company_id").(uint)

	if err := c.stockService.DeleteItem(uint(id), companyID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The product %d has been deleted", id)})
}

func (c *StockController) GetStockByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid stock ID")
	}

	companyID := ctx.Get("company_id").(uint)

	stock, err := c.stockService.GetItemByID(uint(id), companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, stock)
}

func (c *StockController) GetAllItems(ctx echo.Context) error {
	companyID := ctx.Get("company_id").(uint)

	items, err := c.stockService.GetAllItems(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, &items)
}

func (c *StockController) GetDashboard(ctx echo.Context) error {
	companyID := ctx.Get("company_id").(uint)

	dashboard, err := c.stockService.GetDashboard(companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dashboard)
}
