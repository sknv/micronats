package xnats

import (
	"context"
	"fmt"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
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
	argsMsg := &Message{Body: body}

	replyMsg := new(Message)
	if err = c.EncConn.RequestWithContext(ctx, proc, argsMsg, replyMsg); err != nil { // handle network errors
		if err == context.DeadlineExceeded { // handle timeout error if such exist
			err = ErrorStatus(StatusDeadlineExceeded, err.Error())
		}
		return errors.WithMessage(err, fmt.Sprintf("failed to call %s", proc))
	}

	// handle errors transferred over the network
	stat := replyMsg.Status
	if stat.HasError() {
		return stat
	}

	// decode the reply if we are ok
	if err = c.EncConn.Enc.Decode(proc, replyMsg.Body, reply); err != nil {
		return errors.WithMessage(err, "failed to decode the message body")
	}
	return nil
}
