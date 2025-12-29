package controller

import (
	"fmt"
	"net/http"
	"strconv"

	services "github.com/francotraversa/Sliceflow/internal/services/material"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateMaterialHandler godoc
// @Summary      Crear un nuevo material
// @Description  Registra un insumo (Filamento, Resina, etc) para usar en las órdenes de producción.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateMaterialDTO  true  "Datos del Material"
// @Success      201     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Router       /hornero/loged/materials/addmat [post]
func CreateMaterialHandler(c echo.Context) error {
	var newmaterial types.CreateMaterialDTO

	if err := c.Bind(&newmaterial); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}
	if err := services.CreateMaterialUseCase(newmaterial); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "The material has been created")
}

// UpdateMaterialHandler godoc
// @Summary      Actualizar un material existente
// @Description  Modifica los datos de un insumo por su ID.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        id      path    int                         true  "ID del Material"
// @Param        request body    types.UpdateMaterialDTO  true  "Nuevos datos"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /hornero/loged/materials/updmat/{id} [put]
func UpdateMaterialHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Param")
	}
	var mat types.UpdateMaterialDTO
	if err := c.Bind(&mat); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}
	if err := services.UpdateMaterialUseCase(id, mat); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("The material %s has been updated", mat.Name))
}

// DeleteMaterialHandler godoc
// @Summary      Eliminar un material (Borrado Lógico)
// @Description  Marca un insumo como eliminado. No lo borra físicamente de la DB.
// @Tags         Production
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del Material"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /hornero/loged/materials/delmat/{id} [delete]
func DeleteMaterialHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Param")
	}
	var mat types.UpdateMaterialDTO
	if err := c.Bind(&mat); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Json")
	}
	if err := services.DeleteMaterialUseCase(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("The material %s has been deleted", mat.Name))
}

// GetMaterialsHandler godoc
// @Summary      Listar todos los materiales activos
// @Description  Obtiene la lista de insumos disponibles (excluye eliminados).
// @Tags         Production
// @Accept       json
// @Produce      json
// @Success      200  {array}   types.Material
// @Failure      500  {object}  map[string]string
// @Router      /hornero/loged/materials/list [get]
func GetMaterialsHandler(c echo.Context) error {
	var filter types.MaterialFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, "Filters don't work")
	}
	materials, err := services.GetAllMaterialsUseCase(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &materials)
}
