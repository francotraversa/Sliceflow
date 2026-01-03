package controller

import (
	"fmt"
	"net/http"
	"strconv"

	services "github.com/francotraversa/Sliceflow/internal/services/machine"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateMachineHandler godoc
// @Summary      Crear nueva impresora 3D
// @Tags         Machines
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateMachineDTO  true  "Datos Máquina"
// @Router       /hornero/authed/machine/addmac [post]
func CreateMachineHandler(c echo.Context) error {
	var newmachine types.CreateMachineDTO
	if err := c.Bind(&newmachine); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	if err := services.CreateMachineUseCase(newmachine); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("Machine %s has been created", newmachine.Name)})
}

// GetMachinesHandler godoc
// @Summary      Listar impresoras
// @Tags         Machines
// @Produce      json
// @Router       /hornero/authed/machine/list [get]
func GetMachinesHandler(c echo.Context) error {
	var filter types.MachineFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters don't work"})
	}
	machines, err := services.GetAllMachinesUseCase(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, machines)
}

// UpdateMachineHandler godoc
// @Summary      Actualizar impresora
// @Tags         Machines
// @Param        id      path    int                     true  "ID Máquina"
// @Param        request body    types.UpdateMachineDTO  true  "Datos Nuevos"
// @Router       /hornero/authed/machine/updmac/{id} [put]
func UpdateMachineHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto types.UpdateMachineDTO

	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	if err := services.UpdateMachineUseCase(id, dto); err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("Machine %d has been updated", id)})
}

// @Summary      Eliminar una impresora (Borrado lógico)
// @Description  Saca la máquina del listado de disponibles.
// @Tags         Machines
// @Param        id      path    int  true  "ID de la Máquina"
// @Success      200     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /hornero/authed/machine/delmac/{id} [delete]
func DeleteMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	if err := services.DeleteMachineUseCase(id); err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("Machine %d has been deleted", id)})
}
