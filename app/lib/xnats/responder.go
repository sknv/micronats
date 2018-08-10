package xnats

import (
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type Responder struct {
	EncConn *nats.EncodedConn
}

func NewResponder(encConn *nats.EncodedConn) *Responder {
	return &Responder{EncConn: encConn}
}

func (r *Responder) Response(replyTo string, reply interface{}, err error) error {
	var msg *Message

	if err != nil { // transfer error if such exist
		stat, _ := StatusFromError(err)
		msg = &Message{Status: stat}
	} else { // transfer the reply
		body, err := r.EncConn.Enc.Encode(replyTo, reply)
		if err != nil {
			return errors.WithMessage(err, "failed to marshal the reply body")
		}
		msg = &Message{Body: body}
	}

	if err = r.EncConn.Publish(replyTo, msg); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}
