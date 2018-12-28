package nats

import (
	"context"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type Unmarshaller interface {
	Unmarshal(data []byte) error
}

type RemoteClient struct {
	Conn *nats.Conn
}

func (c *RemoteClient) Call(ctx context.Context, subject string, req Marshaller, reply Unmarshaller) error {
	// marshal the outcoming message
	body, err := req.Marshal()
	if err != nil {
		return errors.WithMessage(err, "failed to marshal a message body")
	}

	// create a payload and marshal it
	payload := &Payload{
		Body: body,
		Meta: MetadataFrom(ctx),
	}
	payloadData, err := payload.Marshal()
	if err != nil {
		return errors.WithMessage(err, "failed to marshal payload")
	}

	replyMsg, err := c.Conn.RequestWithContext(ctx, subject, payloadData)
	if err != nil { // handle network errors
		if err == context.DeadlineExceeded { // handle timeout error if such exist
			err = StatusError(DeadlineExceeded, err.Error())
		}
		return errors.WithMessagef(err, "failed to call %s", subject)
	}

	replyPayload := new(Payload)
	if err = replyPayload.Unmarshal(replyMsg.Data); err != nil {
		return errors.WithMessage(err, "failed to unmarshal payload")
	}

	// handle errors transferred over the network
	status := replyPayload.Status
	if status.HasError() {
		return status
	}

	// decode the out if we are ok
	if err = reply.Unmarshal(replyPayload.Body); err != nil {
		return errors.WithMessage(err, "failed to unmarshal a message body")
	}
	return nil
}
