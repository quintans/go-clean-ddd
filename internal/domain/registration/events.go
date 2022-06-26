package registration

import "github.com/quintans/go-clean-ddd/internal/domain"

const EventRegistrationCreated = "RegistrationCreated"

type RegistrationCreatedEvent struct {
	ID    string
	Email domain.Email
}

func (e RegistrationCreatedEvent) Kind() string {
	return EventRegistrationCreated
}

const EventEmailVerified = "EmailVerified"

type EmailVerified struct {
	Email domain.Email
}

func (e EmailVerified) Kind() string {
	return EventEmailVerified
}
