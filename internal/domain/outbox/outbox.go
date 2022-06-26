package outbox

import "errors"

type Outbox struct {
	id      int
	kind    string
	payload []byte
}

func NewOutbox(kind string, payload []byte) (Outbox, error) {
	if kind == "" {
		return Outbox{}, errors.New("required field: Outbox.kind")
	}
	return Outbox{
		kind:    kind,
		payload: payload,
	}, nil
}

func RestoreOutbox(id int, kind string, payload []byte) Outbox {
	return Outbox{
		id:      id,
		kind:    kind,
		payload: payload,
	}
}

func (o Outbox) Id() int {
	return o.id
}

func (o Outbox) Kind() string {
	return o.kind
}

func (o Outbox) Payload() []byte {
	return o.payload
}
