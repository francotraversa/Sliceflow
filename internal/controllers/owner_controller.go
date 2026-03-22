package controller

import (
	"net/http"

	services "github.com/francotraversa/Sliceflow/internal/services/company"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// @Summary      Create Company
// @Description  Creates a new company
// @Tags         Owner
// @Accept       json
// @Produce      json
// @Param      company body types.CompanyCreateDTO true "Company object"
// @Success      200  {object}  types.Response
// @Failure      400  {object}  types.Error
// @Failure      500  {object}  types.Error
// @Router       /hornero/authed/owner/newcompany [post]
func CreateCompanyHandler(c echo.Context) error {
	var NewCompany types.CompanyCreateDTO
	if err := c.Bind(&NewCompany); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "invalid request body"})
	}
	err := services.CreateCompanyUseCase(NewCompany)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to create company"})
	}
	return c.JSON(http.StatusOK, types.Response{Message: "company created successfully"})
}

// @Summary      Get All Companies
// @Description  Retrieves all companies
// @Tags         Owner
// @Produce      json
// @Success      200  {object}  []types.Company
// @Failure      500  {object}  types.Error
// @Router       /hornero/authed/owner/allcompany [get]
func GetAllCompanyHandler(c echo.Context) error {
	companies, err := services.GetAllCompaniesUseCase()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to get companies"})
	}
	return c.JSON(http.StatusOK, companies)
}

// @Summary      Delete Company
// @Description  Deletes a company by ID
// @Tags         Owner
// @Produce      json
// @Param      id path string true "Company ID"
// @Success      200  {object}  types.Response
// @Failure      400  {object}  types.Error
// @Failure      500  {object}  types.Error
// @Router       /hornero/authed/owner/deletecompany/{id} [delete]
func DeleteCompanyHandler(c echo.Context) error {
	IdCompany := c.Param("id")
	if IdCompany == "" {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "IdCompany is required"})
	}
	err := services.DeleteCompanyUseCase(IdCompany)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to delete company"})
	}
	return c.JSON(http.StatusOK, types.Response{Message: "company deleted successfully"})
}
