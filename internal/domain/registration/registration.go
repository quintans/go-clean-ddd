package registration

import (
	"context"

	"github.com/google/uuid"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain"
	libent "github.com/quintans/go-clean-ddd/lib/entity"
	"github.com/quintans/go-clean-ddd/lib/eventbus"
)

type Registration struct {
	core libent.Core

	id       string
	email    domain.Email
	verified bool
}

func NewRegistration(ctx context.Context, email domain.Email, policy domain.UniqueEmailPolicy) (Registration, error) {
	if email.IsZero() {
		return Registration{}, faults.New("registration email is undefined")
	}
	ok, err := policy.IsUnique(ctx, email)
	if err != nil {
		return Registration{}, faults.Wrapf(err, "failed to check uniqueness of email on registration")
	}
	if !ok {
		return Registration{}, faults.Errorf("the provided e-mail %s is not unique", email)
	}

	id := uuid.New().String()
	r := Registration{
		id:       id,
		email:    email,
		verified: false,
	}
	r.core.AddEvent(RegistrationCreatedEvent{ID: id, Email: email})
	return r, nil
}

func Hydrate(id string, email domain.Email, verified bool) (Registration, error) {
	if email.IsZero() {
		return Registration{}, faults.New("registration email is undefined")
	}

	r := Registration{
		id:       id,
		email:    email,
		verified: false,
	}
	return r, nil
}

func (r *Registration) Verify() error {
	if r.verified {
		return domain.ErrNoChange
	}
	r.verified = true
	r.core.AddEvent(EmailVerified{Email: r.email})
	return nil
}

func (r Registration) PopEvents() []eventbus.DomainEvent {
	return r.core.PopEvents()
}

func (r Registration) ID() string {
	return r.id
}

func (r Registration) Email() domain.Email {
	return r.email
}

func (r Registration) Verified() bool {
	return r.verified
}
