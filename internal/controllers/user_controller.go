package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/francotraversa/Sliceflow/internal/auth"
	services "github.com/francotraversa/Sliceflow/internal/services/user"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CreateUserHandler godoc
// @Summary      Registrar un nuevo usuario
// @Description  Crea un usuario en la base de datos. Requiere rol de ADMIN.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        user  body      types.UserCreateCreds  true  "Credenciales del usuario"
// @Success      200   {string}  string                 "The User [username] has been created"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/loged/admin/newuser [post]
func CreateUserHandler(c echo.Context) error {
	var UserCreateCreds types.UserCreateCreds
	if err := c.Bind(&UserCreateCreds); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid JSON")
	}
	err := services.CreateUserUseCase(UserCreateCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("The User %s has been created", UserCreateCreds.Username))
}

// UpdateUserHandler godoc
// @Summary      Actualizar usuario
// @Description  Actualiza datos de un usuario. Si se pasa ID en la URL, requiere ser ADMIN. Si no, actualiza al usuario del token.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int                    false "ID del usuario a editar (solo para Admins)"
// @Param        user  body      types.UserUpdateCreds  true  "Datos a actualizar"
// @Success      200   {string}  string                 "The User ID [id] has been updated"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/loged/updmyuser [patch]
// @Router       /hornero/loged/admin/edituser/{id} [patch]
func UpdateUserHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
	}
	claims, ok := token.Claims.(*auth.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to parse custom claims"})
	}

	requesterID := claims.UserId
	requesterRole := claims.Role

	idParam := c.Param("id")
	var targetID uint

	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID format in URL"})
		}
		targetID = uint(id)
	} else {
		targetID = requesterID
	}

	var updateData types.UserUpdateCreds
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
	}

	err := services.UpdateUserUseCase(targetID, requesterID, requesterRole, updateData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("The User ID %d has been updated", targetID))
}

// DeleteUserHandler godoc
// @Summary      Borrado lógico de usuario
// @Description  Cambia el estado del usuario a 'disabled'. Requiere ser el dueño de la cuenta o ADMIN.
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int                    false "ID del usuario a borrar (solo para Admins)"
// @Success      200   {string}  string                 "The UserID [id] has been deleted"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/loged/deletemyuser [delete]
// @Router       /hornero/loged/admin/deleteuser/{id} [delete]
func DeleteUserHandler(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.JwtCustomClaims)

	requesterID := claims.UserId
	requesterRole := claims.Role

	idParam := c.Param("id")
	var targetID uint

	if idParam != "" {
		id, _ := strconv.ParseUint(idParam, 10, 32)
		targetID = uint(id)
	} else {
		targetID = requesterID
	}

	err := services.DeleteUserUseCase(targetID, requesterID, requesterRole)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("The UserID %d has been deleted", targetID))
}

// GetAllUserHandler godoc
// @Summary      Listar todos los usuarios
// @Description  Obtiene una lista de usuarios filtrada por rol o estado. Solo accesible para ADMIN.
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role    query     string  false  "Filtrar por rol (admin/user)"
// @Param        status  query     string  false  "Filtrar por estado (active/disabled)"
// @Success      200     {array}   types.User
// @Failure      400     {string}  string  "Error en la solicitud"
// @Router       /hornero/loged/admin/alluser [get]
func GetAllUserHandler(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.JwtCustomClaims)

	filterRole := c.QueryParam("role")

	users, err := services.GetAllUserUserUseCase(claims.Role, filterRole)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, users)
}
