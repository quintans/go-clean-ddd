package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/usecase"
)

type RegistrationCommand struct {
	Email string `json:"email"`
}

type UpdateCommand struct {
	Id        string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CustomerDTO struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AdminController manages customer
type AdminController struct {
	commander usecase.Commander
	querier   usecase.Querier
}

func NewController(
	commander usecase.Commander,
	querier usecase.Querier,
) AdminController {
	return AdminController{
		commander: commander,
		querier:   querier,
	}
}

// AddRegistration adds a new customer
func (c AdminController) AddRegistration(ctx echo.Context) error {
	var reg RegistrationCommand
	if err := ctx.Bind(&reg); err != nil {
		return err
	}

	cmd := usecase.RegistrationCommand{
		Email: reg.Email,
	}

	u, err := c.commander.Register(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, u)
}

func (c AdminController) UpdateRegistration(ctx echo.Context) error {
	var reg UpdateCommand
	if err := ctx.Bind(&reg); err != nil {
		return err
	}

	cmd := usecase.UpdateCommand{
		Id:        reg.Id,
		FirstName: reg.FirstName,
		LastName:  reg.LastName,
	}

	err := c.commander.Update(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapBad(ctx, err)
	}

	return ctx.NoContent(http.StatusAccepted)
}

// ListRegistrations lists all customers
func (c AdminController) ListRegistrations(ctx echo.Context) error {
	customers, err := c.querier.GetAllCustomers(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	dtos := toCustomerDTOs(customers)

	return ctx.JSON(http.StatusOK, dtos)
}

func wrapBad(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, err.Error())
}

func wrapError(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusInternalServerError, err.Error())
}

func toCustomerDTOs(in []usecase.CustomerDTO) []CustomerDTO {
	out := make([]CustomerDTO, len(in))
	for k, v := range in {
		out[k] = toCustomerDTO(v)
	}
	return out
}

func toCustomerDTO(c usecase.CustomerDTO) CustomerDTO {
	return CustomerDTO{
		Id:        c.Id,
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}
