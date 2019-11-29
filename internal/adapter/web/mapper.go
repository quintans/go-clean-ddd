package web

import (
	"github.com/quintans/go-clean-ddd/internal/domain/model"
	"github.com/quintans/go-clean-ddd/internal/usecase"
)

func CustomersToCustomerDTOs(in []*model.Customer) []CustomerDTO {
	out := make([]CustomerDTO, len(in))
	for k, v := range in {
		out[k] = CustomerToCustomerDTO(v)
	}
	return out
}

func CustomerToCustomerDTO(c *model.Customer) CustomerDTO {
	return CustomerDTO{
		Uuid:  c.Uuid().String(),
		Email: c.Email(),
	}
}

func RegistrationCommandToRegistration(c RegistrationCommand) usecase.Registration {
	return usecase.Registration{
		Email: c.Email,
	}
}
