package email

import (
	"context"
	"log"
)

type FakeEmailServer struct {
}

func (f FakeEmailServer) Confirm(_ context.Context, destination string, id string) error {
	log.Printf("faking sending confirmation link http://localhost:0000/confirm/%s to %s", id, destination)
	return nil
}
