package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xcontext"
	"github.com/sknv/micronats/app/lib/xnats/status"
)

type ProtoClient struct {
	Conn *nats.Conn
}

func NewProtoClient(conn *nats.Conn) *ProtoClient {
	return &ProtoClient{Conn: conn}
}

func (c *ProtoClient) Call(ctx context.Context, proc string, args proto.Message, reply proto.Message) error {
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.Wrap(err, "failed to marshal args to protobuf")
	}

	timeout, _ := xcontext.Timeout(ctx)
	msg, err := c.Conn.Request(proc, data, timeout)
	if err != nil { // handle network errors
		if err != nats.ErrTimeout { // wrap timeout error if such exist
			err = status.Error(status.DeadlineExceeded, err.Error())
		}
		return errors.Wrapf(err, "failed to call a remote proc: %s", proc)
	}

	// todo: handle errors transfered over the network
	//
	// if status.HasError(msg) {
	// 	rerr := new(status.Status)
	// 	if err = proto.Unmarshal(msg.Body, rerr); err != nil {
	// 		return errors.Wrap(err, "failed to unmarshal an error from protobuf")
	// 	}
	// 	return rerr
	// }

	if err = proto.Unmarshal(msg.Data, reply); err != nil {
		return errors.Wrap(err, "failed to unmarshal a reply from protobuf")
	}
	return nil
}
