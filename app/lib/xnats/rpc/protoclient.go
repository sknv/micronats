package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type ProtoClient struct {
	*RemoteClient
}

func NewProtoClient(conn *nats.Conn) *ProtoClient {
	return &ProtoClient{RemoteClient: NewRemoteClient(conn)}
}

func (c *ProtoClient) Call(ctx context.Context, proc string, args proto.Message, reply proto.Message) error {
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal args to protobuf")
	}

	msg, err := c.RemoteClient.Call(ctx, proc, data)
	if err != nil {
		return errors.WithMessage(err, "failed to call a remote proc")
	}

	// todo: handle errors transfered over the network
	//
	// if status.HasError(msg) {
	// 	rerr := new(status.Status)
	// 	if err = proto.Unmarshal(msg.Body, rerr); err != nil {
	// 		return errors.WithMessage(err, "failed to unmarshal an error from protobuf")
	// 	}
	// 	return rerr
	// }

	if err = proto.Unmarshal(msg.Data, reply); err != nil {
		return errors.WithMessage(err, "failed to unmarshal a reply from protobuf")
	}
	return nil
}
