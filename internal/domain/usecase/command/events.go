package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	devent "github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/lib/event"
)

type RegistrationHandler struct {
	port string
}

func NewRegistrationHandler(port string) RegistrationHandler {
	return RegistrationHandler{
		port: port,
	}
}

func (h RegistrationHandler) Accept(e event.DomainEvent) bool {
	return e.Kind() == devent.EventNewRegistration
}

func (h RegistrationHandler) Handle(ctx context.Context, e event.DomainEvent) error {
	switch t := e.(type) {
	case devent.NewRegistration:
		return h.handleNewRegistration(ctx, t)
	default:
		return faults.Errorf("RegistrationHandler cannot handle event of type %T", e)
	}
}

func (h RegistrationHandler) handleNewRegistration(ctx context.Context, e devent.NewRegistration) error {
	fmt.Println("===> faking send email")
	go func() {
		time.Sleep(time.Second)
		fmt.Println("===> faking user confirmation")
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", h.port))
		if err != nil {
			fmt.Println("ERROR while calling confirmation:", err)
		}
		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("ERROR while reading body:", err)
			}
			fmt.Println("ERROR: response not OK\n", string(body))
		}
	}()
	return nil
}

type EmailVerifiedHandler struct {
	customerRepository usecase.CustomerRepository
	customerView       usecase.CustomerViewRepository
}

func NewEmailVerifiedHandler(customerRepository usecase.CustomerRepository, customerView usecase.CustomerViewRepository) EmailVerifiedHandler {
	return EmailVerifiedHandler{
		customerRepository: customerRepository,
		customerView:       customerView,
	}
}

func (h EmailVerifiedHandler) Accept(e event.DomainEvent) bool {
	return e.Kind() == devent.EventEmailVerified
}

func (h EmailVerifiedHandler) Handle(ctx context.Context, e event.DomainEvent) error {
	switch t := e.(type) {
	case devent.EmailVerified:
		return h.handleEmailVerified(ctx, t)
	default:
		return faults.Errorf("EmailVerifiedHandler cannot handle event of type %T", e)
	}
}

func (h EmailVerifiedHandler) handleEmailVerified(ctx context.Context, e devent.EmailVerified) error {
	customer, err := entity.NewCustomer(ctx, e.Email, usecase.NewUniquenessPolicy(h.customerView))
	if err != nil {
		return faults.Wrap(err)
	}

	if err := h.customerRepository.Create(ctx, customer); err != nil {
		return faults.Wrap(err)
	}
	return nil
}
