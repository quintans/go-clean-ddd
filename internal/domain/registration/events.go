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
