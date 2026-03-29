package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/material"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateMaterialHandler godoc
// @Summary      Create a new material
// @Description  Registers a supply (Filament, Resin, etc.) for use in production orders.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateMaterialDTO  true  "Material data"
// @Success      201     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Router       /hornero/authed/materials/addmat [post]
func CreateMaterialHandler(c echo.Context) error {
	var newmaterial types.CreateMaterialDTO

	if err := c.Bind(&newmaterial); err != nil {
		slog.Warn("materials: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := services.CreateMaterialUseCase(newmaterial, claims.CompanyId); err != nil {
		slog.Error("materials: creation failed", "name", newmaterial.Name, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: created", "name", newmaterial.Name, "company_id", claims.CompanyId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The material %s has been created", newmaterial.Name)})
}

// UpdateMaterialHandler godoc
// @Summary      Update an existing material
// @Description  Modifies a supply's data by its ID.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        id      path    int                         true  "Material ID"
// @Param        request body    types.UpdateMaterialDTO  true  "Updated data"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /hornero/authed/materials/updmat/{id} [put]
func UpdateMaterialHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("materials: invalid ID param", "param", idParam, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Param"})
	}
	var mat types.UpdateMaterialDTO
	if err := c.Bind(&mat); err != nil {
		slog.Warn("materials: invalid request body", "material_id", id, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := services.UpdateMaterialUseCase(id, mat, claims.CompanyId); err != nil {
		slog.Error("materials: update failed", "material_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: updated", "material_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The material %s has been updated", mat.Name)})
}

// DeleteMaterialHandler godoc
// @Summary      Delete a material (Soft-delete)
// @Description  Marks a supply as deleted. Does not physically remove it from the DB.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Material ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /hornero/authed/materials/delmat/{id} [delete]
func DeleteMaterialHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("materials: invalid ID param", "param", idParam, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}
	var mat types.UpdateMaterialDTO
	if err := c.Bind(&mat); err != nil {
		slog.Warn("materials: invalid request body for delete", "material_id", id, "error", err)
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := services.DeleteMaterialUseCase(id, claims.CompanyId); err != nil {
		slog.Error("materials: deletion failed", "material_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: deleted (soft)", "material_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The material %s has been deleted", mat.Name)})
}

// GetMaterialsHandler godoc
// @Summary      List all active materials
// @Description  Returns the list of available supplies (excludes deleted).
// @Tags         Production
// @Accept       json
// @Produce      json
// @Success      200  {array}   types.Material
// @Failure      500  {object}  map[string]string
// @Router      /hornero/authed/materials/list [get]
func GetMaterialsHandler(c echo.Context) error {
	var filter types.MaterialFilter

	if err := c.Bind(&filter); err != nil {
		slog.Warn("materials: invalid filter params", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters don't work"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	materials, err := services.GetAllMaterialsUseCase(filter, claims.CompanyId)
	if err != nil {
		slog.Error("materials: list failed", "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: listed", "count", len(*materials), "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, &materials)
}
