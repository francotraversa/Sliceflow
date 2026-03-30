package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/user"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// CreateUserHandler godoc
// @Summary      Register a new user
// @Description  Creates a user in the database. Requires ADMIN role.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      types.UserCreateCreds  true  "User credentials"
// @Success      200   {string}  string                 "The User [username] has been created"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/authed/admin/newuser [post]
func CreateUserHandler(c echo.Context) error {
	var UserCreateCreds types.UserCreateCreds
	if err := c.Bind(&UserCreateCreds); err != nil {
		slog.Warn("users: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if claims.Role != "owner" && claims.Role != "superadmin" {
		slog.Warn("users: unauthorized creation attempt", "requester_role", claims.Role, "requester_id", claims.UserId)
		return c.JSON(http.StatusForbidden, types.Error{Error: "Only owners can create new users"})
	}

	err = services.CreateUserUseCase(UserCreateCreds)
	if err != nil {
		slog.Error("users: creation failed", "username", UserCreateCreds.Username, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: created", "username", UserCreateCreds.Username, "by", claims.UserId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The User %s has been created", UserCreateCreds.Username)})
}

// UpdateUserHandler godoc
// @Summary      Update user
// @Description  Updates user data. If ID is in the URL, requires ADMIN role. Otherwise updates the token owner.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        id    path      int                    false "User ID to edit (Admins only)"
// @Param        user  body      types.UserUpdateCreds  true  "Data to update"
// @Success      200   {string}  string                 "The User ID [id] has been updated"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/authed/updmyuser [patch]
// @Router       /hornero/authed/admin/edituser/{id} [patch]
func UpdateUserHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
	}

	requesterID := claims.UserId
	requesterRole := claims.Role

	idParam := c.Param("id")
	var targetID uint

	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			slog.Warn("users: invalid ID param", "param", idParam, "error", err)
			return c.JSON(http.StatusInternalServerError, types.Error{Error: "invalid user ID format in URL"})
		}
		targetID = uint(id)
	} else {
		targetID = requesterID
	}

	var updateData types.UserUpdateCreds
	if err := c.Bind(&updateData); err != nil {
		slog.Warn("users: invalid request body", "target_id", targetID, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "invalid JSON body"})
	}

	err = services.UpdateUserUseCase(targetID, requesterID, requesterRole, updateData)
	if err != nil {
		slog.Error("users: update failed", "target_id", targetID, "requester_id", requesterID, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: updated", "target_id", targetID, "by", requesterID)
	return c.JSON(http.StatusAccepted, types.Response{Message: fmt.Sprintf("The User %s has been updated", updateData.Username)})
}

// DeleteUserHandler godoc
// @Summary      Soft-delete user
// @Description  Sets the user status to 'disabled'. Requires account ownership or ADMIN role.
// @Tags         Users
// @Produce      json
// @Security BearerAuth
// @Param        id    path      int                    false "User ID to delete (Admins only)"
// @Success      200   {string}  string                 "The UserID [id] has been deleted"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/authed/deletemyuser [delete]
// @Router       /hornero/authed/admin/deleteuser/{id} [delete]
func DeleteUserHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
	}

	requesterID := claims.UserId
	requesterRole := claims.Role

	idParam := c.Param("id")
	var targetID uint

	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			slog.Warn("users: invalid ID param", "param", idParam, "error", err)
			return c.JSON(http.StatusInternalServerError, types.Error{Error: "invalid user ID format in URL"})
		}
		targetID = uint(id)
	} else {
		targetID = requesterID
	}

	err = services.DeleteUserUseCase(targetID, requesterID, requesterRole)
	if err != nil {
		slog.Error("users: deletion failed", "target_id", targetID, "requester_id", requesterID, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: deleted (soft)", "target_id", targetID, "by", requesterID)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The UserID %d has been deleted", targetID)})
}

// EnableUserHandler godoc
// @Summary      Enable user by ID
// @Description  Sets the user status to 'active'. Requires ADMIN permissions.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      types.UserIDRequest  true  "JSON with the user ID"
// @Success      200   {object}  types.Response       "User with ID [id] has been enabled"
// @Failure      400   {object}  types.Response       "Error message"
// @Router       /hornero/authed/admin/enableuser [delete]
func EnableUserHandler(c echo.Context) error {
	var req types.UserIDActivate

	if err := c.Bind(&req); err != nil {
		slog.Warn("users: invalid request body for enable", "error", err)
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: "Invalid input format",
		})
	}

	err := services.EnableUserByIDUseCase(req)
	if err != nil {
		slog.Error("users: enable failed", "user_id", req.IdUser, "error", err)
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: err.Error(),
		})
	}

	slog.Info("users: enabled", "user_id", req.IdUser)
	return c.JSON(http.StatusOK, types.Response{
		Message: fmt.Sprintf("User with ID %d has been enabled successfully", req.IdUser),
	})
}

// GetAllUserHandler godoc
// @Summary      List all users
// @Description  Returns a list of users filtered by role or status. Admin only.
// @Tags         Users
// @Produce      json
// @Security BearerAuth
// @Param        role    query     string  false  "Filter by role (admin/user)"
// @Param        status  query     string  false  "Filter by status (active/disabled)"
// @Success      200     {array}   types.User
// @Failure      400     {string}  string  "Error message"
// @Router       /hornero/authed/admin/alluser [get]
func GetAllUserHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
	}

	filterRole := c.QueryParam("role")

	users, err := services.GetAllUserUserUseCase(claims.Role, filterRole, claims.Role, int(claims.CompanyId))
	if err != nil {
		slog.Error("users: list failed", "company_id", claims.CompanyId, "filter_role", filterRole, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: listed", "count", len(*users), "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, users)
}

// CreateAdminHandler godoc
// @Summary      Create new admin
// @Description  Creates a new admin user. Requires owner permissions.
// @Tags         Users
// @Produce      json
// @Security BearerAuth
// @Param        request  body      types.UserCreateCreds  true  "Admin credentials"
// @Success      200   {object}  types.Response       "Admin created successfully"
// @Failure      400   {object}  types.Response       "Error message"
// @Router       /hornero/authed/owner/newadmin [post]
func CreateAdminHandler(c echo.Context) error {
	var UserCreateCreds types.UserCreateCreds
	if err := c.Bind(&UserCreateCreds); err != nil {
		slog.Warn("users: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if claims.Role != "owner" {
		slog.Warn("users: unauthorized creation attempt", "requester_role", claims.Role, "requester_id", claims.UserId)
		return c.JSON(http.StatusForbidden, types.Error{Error: "Only owners can create new users"})
	}

	err = services.CreateAdminUseCase(UserCreateCreds)
	if err != nil {
		slog.Error("users: creation failed", "username", UserCreateCreds.Username, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: created", "username", UserCreateCreds.Username, "by", claims.UserId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The User %s has been created", UserCreateCreds.Username)})
}

// DeleteAdminHandler godoc
// @Summary      Delete admin by ID
// @Description  Sets the user status to 'disabled'. Requires owner permissions.
// @Tags         Users
// @Produce      json
// @Security BearerAuth
// @Param        id    path      int                    false "User ID to delete (Admins only)"
// @Success      200   {string}  string                 "The UserID [id] has been deleted"
// @Failure      400   {string}  string                 "Error message"
// @Router       /hornero/authed/owner/deleteadmin/{id} [delete]
func DeleteAdminHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
	}

	idParam := c.Param("id")
	var targetID uint

	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			slog.Warn("users: invalid ID param", "param", idParam, "error", err)
			return c.JSON(http.StatusInternalServerError, types.Error{Error: "invalid user ID format in URL"})
		}
		targetID = uint(id)
	} else {
		targetID = claims.UserId
	}

	err = services.DeleteAdminUseCase(targetID)
	if err != nil {
		slog.Error("users: deletion failed", "target_id", targetID, "requester_id", claims.UserId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: deleted (soft)", "target_id", targetID, "by", claims.UserId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The UserID %d has been deleted", targetID)})
}

func GetAllAdminHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("users: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
	}
	if claims.Role != "owner" {
		slog.Warn("users: unauthorized list attempt", "requester_role", claims.Role, "requester_id", claims.UserId)
		return c.JSON(http.StatusForbidden, types.Error{Error: "Only owners can list admins"})
	}

	users, err := services.GetAllAdminUserUseCase()
	if err != nil {
		slog.Error("users: list failed", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("users: listed", "count", len(*users))
	return c.JSON(http.StatusOK, users)
}
