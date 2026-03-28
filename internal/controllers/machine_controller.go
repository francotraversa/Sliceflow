package controller

import (
	"fmt"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/machine"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateMachineHandler godoc
// @Summary      Create a new 3D printer
// @Tags         Machines
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateMachineDTO  true  "Machine data"
// @Router       /hornero/authed/machine/addmac [post]
func CreateMachineHandler(c echo.Context) error {
	var newmachine types.CreateMachineDTO
	if err := c.Bind(&newmachine); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := services.CreateMachineUseCase(newmachine, claims.CompanyId); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("Machine %s has been created", newmachine.Name)})
}

// GetMachinesHandler godoc
// @Summary      List printers
// @Tags         Machines
// @Produce      json
// @Router       /hornero/authed/machine/list [get]
func GetMachinesHandler(c echo.Context) error {
	var filter types.MachineFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters don't work"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	machines, err := services.GetAllMachinesUseCase(filter, claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, machines)
}

// UpdateMachineHandler godoc
// @Summary      Update printer
// @Tags         Machines
// @Param        id      path    int                     true  "Machine ID"
// @Param        request body    types.UpdateMachineDTO  true  "Updated data"
// @Router       /hornero/authed/machine/updmac/{id} [put]
func UpdateMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "invalid machine ID format in URL"})
	}
	var dto types.UpdateMachineDTO

	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := services.UpdateMachineUseCase(id, dto, claims.CompanyId); err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Machine %d has been updated", id)})
}

// @Summary      Delete a printer (Soft-delete)
// @Description  Removes the machine from the available list.
// @Tags         Machines
// @Param        id      path    int  true  "Machine ID"
// @Success      200     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /hornero/authed/machine/delmac/{id} [delete]
func DeleteMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := services.DeleteMachineUseCase(id, claims.CompanyId); err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("Machine %d has been deleted", id)})
}
