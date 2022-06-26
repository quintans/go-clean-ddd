package app

import (
	"context"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain"
)

type UniquenessPolicer interface {
	IsUnique(ctx context.Context, email domain.Email) (bool, error)
}

type UniquenessPolicy struct {
	customerView CustomerViewRepository
}

func NewUniquenessPolicy(customerView CustomerViewRepository) UniquenessPolicy {
	return UniquenessPolicy{
		customerView: customerView,
	}
}

func (p UniquenessPolicy) IsUnique(ctx context.Context, email domain.Email) (bool, error) {
	_, err := p.customerView.GetByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return err == nil, faults.Wrap(err)
}
