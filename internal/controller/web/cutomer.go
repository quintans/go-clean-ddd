package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/usecase/command"
	"github.com/quintans/go-clean-ddd/internal/usecase/query"
)

type UpdateCommand struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CustomerDTO struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// CustomerController manages customer
type CustomerController struct {
	updateCustomerHandler command.UpdateCustomerHandler
	allCustomersHandler   query.AllCustomersHandler
}

func NewCustomerController(
	updateCustomerHandler command.UpdateCustomerHandler,
	allCustomersHandler query.AllCustomersHandler,
) CustomerController {
	return CustomerController{
		updateCustomerHandler: updateCustomerHandler,
		allCustomersHandler:   allCustomersHandler,
	}
}

func (c CustomerController) UpdateRegistration(ctx echo.Context) error {
	id := ctx.Param("id")

	var reg UpdateCommand
	if err := ctx.Bind(&reg); err != nil {
		return err
	}

	cmd := command.UpdateCustomerCommand{
		Id:        id,
		FirstName: reg.FirstName,
		LastName:  reg.LastName,
	}

	err := c.updateCustomerHandler.Handle(ctx.Request().Context(), cmd)
	if err != nil {
		return wrapBad(ctx, err)
	}

	return ctx.NoContent(http.StatusAccepted)
}

// ListRegistrations lists all customers
func (c CustomerController) ListRegistrations(ctx echo.Context) error {
	customers, err := c.allCustomersHandler.Handle(ctx.Request().Context())
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

func toCustomerDTOs(in []query.CustomerDTO) []CustomerDTO {
	out := make([]CustomerDTO, len(in))
	for k, v := range in {
		out[k] = toCustomerDTO(v)
	}
	return out
}

func toCustomerDTO(c query.CustomerDTO) CustomerDTO {
	return CustomerDTO{
		Id:        c.Id,
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}
