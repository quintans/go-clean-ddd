package persistence

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/model"
	"github.com/quintans/go-clean-ddd/internal/usecase"
)

var _ usecase.RegistrationService = (*RegistrationServiceTx)(nil)

func NewRegistrationServiceTx(tm *TxManager, service usecase.RegistrationService) usecase.RegistrationService {
	return RegistrationServiceTx{tm, service}
}

type RegistrationServiceTx struct {
	tm   *TxManager
	next usecase.RegistrationService
}

func (s RegistrationServiceTx) FindAllResgistrations(ctx context.Context) (customers []*model.Customer, err error) {
	s.tm.None(ctx).Do(func(ctx context.Context) error {
		customers, err = s.next.FindAllResgistrations(ctx)
		return err
	})
	return
}

func (s RegistrationServiceTx) Register(ctx context.Context, registerDto usecase.Registration) (customer *model.Customer, err error) {
	s.tm.Current(ctx).Do(func(ctx context.Context) error {
		customer, err = s.next.Register(ctx, registerDto)
		return err
	})
	return
}
