package web

import (
	"net/http"

	"github.com/quintans/go-clean-ddd/internal/usecase"
	"github.com/labstack/echo/v4"
)

type RegistrationCommand struct {
	Email string `json:"email"`
}

type CustomerDTO struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
}

// AdminController manages customer
type AdminController struct {
	registrationService usecase.RegistrationService
}

func NewAdminController(registrationService usecase.RegistrationService) AdminController {
	return AdminController{
		registrationService: registrationService,
	}
}

// AddRegistration adds a new customer
func (c AdminController) AddRegistration(ctx echo.Context) error {
	var dto RegistrationCommand
	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	u, err := c.registrationService.Register(ctx.Request().Context(), RegistrationCommandToRegistration(dto))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, u)
}

// ListRegistrations lists all customers
func (c AdminController) ListRegistrations(ctx echo.Context) error {
	customers, err := c.registrationService.FindAllResgistrations(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	dtos := CustomersToCustomerDTOs(customers)

	return ctx.JSON(http.StatusOK, dtos)
}
