package xnats

import (
	"context"
	"fmt"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats/message"
	"github.com/sknv/micronats/app/lib/xnats/status"
)

type RemoteClient struct {
	EncConn *nats.EncodedConn
}

func NewRemoteClient(encConn *nats.EncodedConn) *RemoteClient {
	return &RemoteClient{EncConn: encConn}
}

func (c *RemoteClient) Call(ctx context.Context, proc string, args interface{}, reply interface{}) error {
	// encode the args message
	body, err := c.EncConn.Enc.Encode(proc, args)
	if err != nil {
		return errors.WithMessage(err, "failed to encode the message body")
	}

	// todo: add metadata from context
	argsMsg := &message.Message{Body: body}

	replyMsg := new(message.Message)
	if err = c.EncConn.RequestWithContext(ctx, proc, argsMsg, replyMsg); err != nil { // handle network errors
		if err == context.DeadlineExceeded { // handle timeout error if such exist
			err = status.Error(status.DeadlineExceeded, err.Error())
		}
		return errors.WithMessage(err, fmt.Sprintf("failed to call %s", proc))
	}

	// handle errors transferred over the network
	status := replyMsg.Status
	if status.HasError() {
		return status
	}

	// decode the reply if we are ok
	if err = c.EncConn.Enc.Decode(proc, replyMsg.Body, reply); err != nil {
		return errors.WithMessage(err, "failed to decode the message body")
	}
	return nil
}
