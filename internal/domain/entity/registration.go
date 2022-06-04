package entity

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

type Registration struct {
	Core

	id       string
	email    vo.Email
	verified bool
}

func NewRegistration(ctx context.Context, email vo.Email, policy UniqueEmailPolicy) (Registration, error) {
	if email.IsZero() {
		return Registration{}, errors.New("registration email is undefined")
	}
	ok, err := policy.IsUnique(ctx, email)
	if err != nil {
		return Registration{}, errors.Wrap(err, "failed to check uniqueness of email")
	}
	if !ok {
		return Registration{}, errors.Errorf("the provided e-mail %s is not unique", email)
	}

	id := uuid.New().String()
	r := Registration{
		id:       id,
		email:    email,
		verified: false,
	}
	r.AddEvent(event.NewRegistration{Id: id})
	return r, nil
}

func RestoreRegistration(id string, email vo.Email, verified bool) (Registration, error) {
	if email.IsZero() {
		return Registration{}, errors.New("registration email is undefined")
	}

	r := Registration{
		id:       id,
		email:    email,
		verified: false,
	}
	return r, nil
}

func (r *Registration) Verify() {
	if r.verified {
		return
	}
	r.verified = true
	r.AddEvent(event.EmailVerified{Email: r.email})
}

func (r Registration) ID() string {
	return r.id
}

func (r Registration) Email() vo.Email {
	return r.email
}

func (r Registration) Verified() bool {
	return r.verified
}
