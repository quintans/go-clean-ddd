package fakemq

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
)

type FakePublisher struct{}

func (f FakePublisher) Publish(_ context.Context, event usecase.NewRegistration) error {
	url := "http://localhost:0000/confirm/" + event.Id
	log.Printf("faking sending confirmation link %s to %s\n", url, event.Email)
	go func() {
		time.Sleep(500 * time.Millisecond)
		_, err := http.Get(url)
		if err != nil {
			log.Printf("ERROR: failed to call %s: %s\n", url, err)
		}
	}()
	return nil
}
