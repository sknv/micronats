package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats/message"
	"github.com/sknv/micronats/app/lib/xnats/status"
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

	// unmarshal to protowrapper
	protoMsg := new(message.Message)
	if err = proto.Unmarshal(msg.Data, protoMsg); err != nil {
		return errors.WithMessage(err, "failed to unmarshal the reply from protobuf")
	}

	// handle errors transfered over the network
	if protoMsg.HasError() {
		status := new(status.Status)
		if err = ptypes.UnmarshalAny(protoMsg.Body, status); err != nil {
			return errors.WithMessage(err, "failed to unmarshal the error from protobug any")
		}
		return status
	}

	if err = ptypes.UnmarshalAny(protoMsg.Body, reply); err != nil {
		return errors.WithMessage(err, "failed to unmarshal the reply from protobuf any")
	}
	return nil
}
