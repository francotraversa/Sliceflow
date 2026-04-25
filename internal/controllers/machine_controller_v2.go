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

type MachineController struct {
	machineUseCase domain.MachineUseCase
}

func NewMachineController(uc domain.MachineUseCase) *MachineController {
	return &MachineController{machineUseCase: uc}
}

func (mc *MachineController) CreateMachineHandler(c echo.Context) error {
	var newmachine types.CreateMachineDTO
	if err := c.Bind(&newmachine); err != nil {
		slog.Warn("machines: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("machines: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	err = mc.machineUseCase.CreateMachine(newmachine, claims.CompanyId)
	if err != nil {
		slog.Error("machines: creation failed", "name", newmachine.Name, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("machines: created", "name", newmachine.Name, "company_id", claims.CompanyId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("Machine %s has been created", newmachine.Name)})
}

func (mc *MachineController) GetMachinesHandler(c echo.Context) error {
	var filter types.MachineFilter
	if err := c.Bind(&filter); err != nil {
		slog.Warn("machines: invalid filter params", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters don't work"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("machines: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	machines, err := mc.machineUseCase.GetMachines(filter, claims.CompanyId)
	if err != nil {
		slog.Error("machines: list failed", "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("machines: listed", "count", len(machines), "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, machines)
}

func (mc *MachineController) GetMachineByIDHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("machines: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("machines: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	machine, err := mc.machineUseCase.GetMachineByID(uint(id), claims.CompanyId)
	if err != nil {
		slog.Error("machines: get by ID failed", "machine_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("machines: got by ID", "machine_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, machine)
}

func (mc *MachineController) UpdateMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("machines: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "invalid machine ID format in URL"})
	}
	var updatedMachine types.UpdateMachineDTO
	if err := c.Bind(&updatedMachine); err != nil {
		slog.Warn("machines: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("machines: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	err = mc.machineUseCase.UpdateMachine(uint(id), updatedMachine, claims.CompanyId)
	if err != nil {
		slog.Error("machines: update failed", "name", updatedMachine.Name, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("machines: updated", "name", updatedMachine.Name, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Machine %d has been updated", id)})
}

func (mc *MachineController) DeleteMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("machines: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Id"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("machines: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	err = mc.machineUseCase.DeleteMachine(uint(id), claims.CompanyId)
	if err != nil {
		slog.Error("machines: delete failed", "machine_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	slog.Info("machines: deleted", "machine_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Machine %d has been deleted", id)})
}
