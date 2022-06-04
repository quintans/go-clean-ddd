package event

import "github.com/quintans/go-clean-ddd/internal/domain/vo"

const EventEmailVerified = "EmailVerified"

type EmailVerified struct {
	Email vo.Email
}

func (e EmailVerified) Kind() string {
	return EventEmailVerified
}

const EventNewRegistration = "NewRegistration"

type NewRegistration struct {
	Id    string
	Email vo.Email
}

func (e NewRegistration) Kind() string {
	return EventNewRegistration
}
