package nats

import (
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type Msg struct {
	*nats.Msg

	payload *Payload
}

func (m *Msg) Payload() (*Payload, error) {
	if m.payload != nil {
		return m.payload, nil
	}

	out := new(Payload)
	if err := out.Unmarshal(m.Data); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal payload")
	}
	m.payload = out
	return out, nil
}
