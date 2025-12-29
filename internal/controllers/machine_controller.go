package controller

import (
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
// @Router       /hornero/loged/machine/addmac [post]
func CreateMachineHandler(c echo.Context) error {
	var newmachine types.CreateMachineDTO
	if err := c.Bind(&newmachine); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}

	if err := services.CreateMachineUseCase(newmachine); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, "Machine has been created")
}

// GetMachinesHandler godoc
// @Summary      Listar impresoras
// @Tags         Machines
// @Produce      json
// @Router       /hornero/loged/machine/list [get]
func GetMachinesHandler(c echo.Context) error {
	var filter types.MachineFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, "Filter don't work")
	}
	machines, err := services.GetAllMachinesUseCase(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, machines)
}

// UpdateMachineHandler godoc
// @Summary      Actualizar impresora
// @Tags         Machines
// @Param        id      path    int                     true  "ID Máquina"
// @Param        request body    types.UpdateMachineDTO  true  "Datos Nuevos"
// @Router       /hornero/loged/machine/updmac/{id} [put]
func UpdateMachineHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto types.UpdateMachineDTO

	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}

	if err := services.UpdateMachineUseCase(id, dto); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "The machine has been updated")
}

// @Summary      Eliminar una impresora (Borrado lógico)
// @Description  Saca la máquina del listado de disponibles.
// @Tags         Machines
// @Param        id      path    int  true  "ID de la Máquina"
// @Success      200     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /hornero/loged/machine/delmac/{id} [delete]
func DeleteMachineHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID Param")
	}

	if err := services.DeleteMachineUseCase(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "The machine has been disabled")
}
