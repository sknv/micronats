package rpc

import (
	"context"
	"fmt"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xcontext"
	"github.com/sknv/micronats/app/lib/xnats/status"
)

type RemoteClient struct {
	Conn *nats.Conn
}

func NewRemoteClient(conn *nats.Conn) *RemoteClient {
	return &RemoteClient{Conn: conn}
}

func (c *RemoteClient) Call(ctx context.Context, proc string, data []byte) (*nats.Msg, error) {
	timeout, _ := xcontext.Timeout(ctx)
	msg, err := c.Conn.Request(proc, data, timeout)
	if err != nil { // handle network errors
		if err == nats.ErrTimeout { // wrap timeout error if such exist
			err = status.Error(status.DeadlineExceeded, err.Error())
		}
		return nil, errors.WithMessage(err, fmt.Sprintf("failed to make a request to %s", proc))
	}
	return msg, nil
}
