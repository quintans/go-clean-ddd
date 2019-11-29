package usecase

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/model"
	"github.com/quintans/go-clean-ddd/internal/domain/service"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

/*
The application services are responsible for driving workflow and coordinating transaction management.
They also provide a high-level abstraction for clients to use when interacting with the domain.
These services are typically designed to define or support specific use cases.
*/

type RegistrationService interface {
	FindAllResgistrations(context.Context) ([]*model.Customer, error)
	Register(context.Context, Registration) (*model.Customer, error)
}

type Registration struct {
	Email string
}

type RegistrationServiceImpl struct {
	customerService    service.CustomerService
	customerRepository model.CustomerRepository
}

func NewRegistrationServiceImpl(customerService service.CustomerService, customerRepository model.CustomerRepository) RegistrationServiceImpl {
	return RegistrationServiceImpl{
		customerService:    customerService,
		customerRepository: customerRepository,
	}
}

func (r RegistrationServiceImpl) FindAllResgistrations(ctx context.Context) ([]*model.Customer, error) {
	return r.customerRepository.FindAll(ctx)
}

func (r RegistrationServiceImpl) Register(ctx context.Context, reg Registration) (*model.Customer, error) {
	old, err := r.customerService.FindDuplicate(ctx, reg.Email)
	if err != nil {
		return nil, err
	}
	if old == nil {
		customer := model.NewCustomerAll(
			uuid.New(),
			model.NewFullNameVO("", ""),
			reg.Email,
		)

		if err = customer.SetEmail(reg.Email); err != nil {
			return nil, err
		}

		err = r.customerRepository.Store(ctx, customer)
		return customer, err
	}

	return nil, errors.New("Customer with the same email already exists")
}
