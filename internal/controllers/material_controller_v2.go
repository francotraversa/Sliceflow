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

type MaterialController struct {
	materialUseCase domain.MaterialUseCase
}

func NewMaterialController(uc domain.MaterialUseCase) *MaterialController {
	return &MaterialController{materialUseCase: uc}
}

func (mc *MaterialController) CreateMaterialHandler(c echo.Context) error {
	var newMaterial types.CreateMaterialDTO
	if err := c.Bind(&newMaterial); err != nil {
		slog.Warn("materials: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	err = mc.materialUseCase.CreateMaterial(newMaterial, claims.CompanyId)
	if err != nil {
		slog.Error("materials: creation failed", "name", newMaterial.Name, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: created", "name", newMaterial.Name, "company_id", claims.CompanyId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("Material %s has been created", newMaterial.Name)})
}

func (mc *MaterialController) GetMaterialsHandler(c echo.Context) error {
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
	materials, err := mc.materialUseCase.GetMaterials(filter, claims.CompanyId)
	if err != nil {
		slog.Error("materials: list failed", "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("materials: listed", "count", len(materials), "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, materials)
}

func (mc *MaterialController) GetMaterialByIDHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("materials: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	material, err := mc.materialUseCase.GetMaterialByID(uint(id), claims.CompanyId)
	if err != nil {
		slog.Error("materials: get by ID failed", "material_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("materials: got by ID", "material_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, material)
}

func (mc *MaterialController) UpdateMaterialHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("materials: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}
	var updatedMaterial types.UpdateMaterialDTO
	if err := c.Bind(&updatedMaterial); err != nil {
		slog.Warn("materials: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	err = mc.materialUseCase.UpdateMaterial(uint(id), updatedMaterial, claims.CompanyId)
	if err != nil {
		slog.Error("materials: update failed", "name", updatedMaterial.Name, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("materials: updated", "name", updatedMaterial.Name, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Material %d has been updated", id)})
}

func (mc *MaterialController) DeleteMaterialHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("materials: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("materials: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	err = mc.materialUseCase.DeleteMaterial(uint(id), claims.CompanyId)
	if err != nil {
		slog.Error("materials: delete failed", "material_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("materials: deleted", "material_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Material %d has been deleted", id)})
}
