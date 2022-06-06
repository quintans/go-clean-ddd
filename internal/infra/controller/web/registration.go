package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
)

// RegistrationController manages customer
type RegistrationController struct {
	createRegistrationHandler command.CreateRegistrationHandler
}

func NewRegistrationController(
	createRegistrationHandler command.CreateRegistrationHandler,
) RegistrationController {
	return RegistrationController{
		createRegistrationHandler: createRegistrationHandler,
	}
}

type RegistrationCommand struct {
	Email string `json:"email"`
}

// AddRegistration adds a new customer
func (c RegistrationController) AddRegistration(ctx echo.Context) error {
	var reg RegistrationCommand
	if err := ctx.Bind(&reg); err != nil {
		return err
	}

	cmd := command.CreateRegistrationCommand{
		Email: reg.Email,
	}

	u, err := c.createRegistrationHandler.Handle(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, u)
}
