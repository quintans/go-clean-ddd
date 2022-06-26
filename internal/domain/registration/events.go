package registration

import "github.com/quintans/go-clean-ddd/internal/domain"

const EventRegistered = "Registered"

type RegisteredEvent struct {
	Id    string
	Email domain.Email
}

func (e RegisteredEvent) Kind() string {
	return EventRegistered
}

const EventEmailVerified = "EmailVerified"

type EmailVerified struct {
	Email domain.Email
}

func (e EmailVerified) Kind() string {
	return EventEmailVerified
}
