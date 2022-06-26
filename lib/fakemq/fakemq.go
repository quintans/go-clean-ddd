package fakemq

import (
	"context"
	"sync"
)

type FakeMQ struct {
	channel     chan Event
	subscribers map[string][]Subscriber

	mu    sync.RWMutex
	close bool
}

type Event struct {
	Kind    string
	Payload []byte
}

type Subscriber interface {
	Handle(context.Context, Event) error
}

func New() *FakeMQ {
	return &FakeMQ{
		channel: make(chan Event, 10),
	}
}

func (f *FakeMQ) Publish(event Event) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.close {
		return
	}

	f.channel <- event
}

func (f *FakeMQ) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.close {
		return
	}
	close(f.channel)
	f.close = true
}

func (f *FakeMQ) Subscribe(kind string, s Subscriber) {
	subs := f.subscribers[kind]
	subs = append(subs, s)
	f.subscribers[kind] = subs
}

func (f *FakeMQ) Start() {
	go func() {
		for e := range f.channel {
			subs := f.subscribers[e.Kind]
			for _, s := range subs {
				s.Handle(context.Background(), e)
			}
		}
	}()
}
