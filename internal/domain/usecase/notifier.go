package usecase

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

type Publisher interface {
	Publish(ctx context.Context, event NewRegistration) error
}

type NewRegistration struct {
	Id    string
	Email vo.Email
}
