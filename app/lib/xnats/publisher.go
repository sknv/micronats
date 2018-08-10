package xnats

import (
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type Publisher struct {
	EncConn *nats.EncodedConn
}

func NewPublisher(encConn *nats.EncodedConn) *Publisher {
	return &Publisher{EncConn: encConn}
}

func (r *Publisher) Publish(subject string, message interface{}, err error) error {
	var msg *Message
	// prepare the reply
	if err != nil { // transfer error if such exist
		stat, _ := StatusFromError(err)
		msg = &Message{Status: stat}
	} else { // transfer the reply
		body, err := r.EncConn.Enc.Encode(subject, message)
		if err != nil {
			return errors.WithMessage(err, "failed to marshal the reply body")
		}
		msg = &Message{Body: body}
	}

	if err = r.EncConn.Publish(subject, msg); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}
