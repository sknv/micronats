package nats

import (
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type Marshaller interface {
	Marshal() ([]byte, error)
}

type Responder struct {
	Conn *nats.Conn
}

func (r *Responder) Respond(subject string, response Marshaller, err error) error {
	payload := new(Payload)
	if err != nil { // transfer an error if such exist
		status, _ := StatusFromError(err)
		payload.Status = status
	} else { // otherwise, transfer the response
		body, err := response.Marshal()
		if err != nil {
			return errors.WithMessage(err, "failed to marshal a response")
		}
		payload.Body = body
	}

	// marshal payload
	outData, err := payload.Marshal()
	if err != nil {
		return errors.WithMessage(err, "failed to marshal payload")
	}

	if err = r.Conn.Publish(subject, outData); err != nil {
		return errors.WithMessage(err, "failed to publish a message")
	}
	return nil
}
