package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app/command"
)

// RegistrationController manages customer
type RegistrationController struct {
	createRegistrationHandler  command.CreateRegistrationHandler
	confirmRegistrationHandler command.ConfirmRegistrationHandler
}

func NewRegistrationController(
	createRegistrationHandler command.CreateRegistrationHandler,
	confirmRegistrationHandler command.ConfirmRegistrationHandler,
) RegistrationController {
	return RegistrationController{
		createRegistrationHandler:  createRegistrationHandler,
		confirmRegistrationHandler: confirmRegistrationHandler,
	}
}

type RegistrationCommand struct {
	Email string `json:"email"`
}

// AddRegistration adds a new customer
func (c RegistrationController) AddRegistration(ctx echo.Context) error {
	var reg RegistrationCommand
	if err := ctx.Bind(&reg); err != nil {
		return faults.Wrap(err)
	}

	cmd := command.CreateRegistrationCommand{
		Email: reg.Email,
	}

	u, err := c.createRegistrationHandler.Handle(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, u)
}

func (c RegistrationController) ConfirmRegistration(ctx echo.Context) error {
	id := ctx.Param("id")

	cmd := command.ConfirmRegistrationCommand{
		Id: id,
	}

	err := c.confirmRegistrationHandler.Handle(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapError(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}
