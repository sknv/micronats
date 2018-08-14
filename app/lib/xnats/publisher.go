package xnats

import (
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats/message"
	"github.com/sknv/micronats/app/lib/xnats/status"
)

type Publisher struct {
	EncConn *nats.EncodedConn
}

func NewPublisher(encConn *nats.EncodedConn) *Publisher {
	return &Publisher{EncConn: encConn}
}

func (r *Publisher) Publish(subject string, publishing interface{}, err error) error {
	var msg *message.Message
	if err != nil { // transfer error if such exist
		status, _ := status.FromError(err)
		msg = &message.Message{Status: status}
	} else { // transfer the publishing
		body, err := r.EncConn.Enc.Encode(subject, publishing)
		if err != nil {
			return errors.WithMessage(err, "failed to marshal the publishing body")
		}
		msg = &message.Message{Body: body}
	}

	if err = r.EncConn.Publish(subject, msg); err != nil {
		return errors.WithMessage(err, "failed to publish the message")
	}
	return nil
}
