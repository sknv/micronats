package nats

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type ProtoClient struct {
	Conn    *nats.Conn
	Timeout time.Duration
}

func NewProtoClient(natsconn *nats.Conn, timeout time.Duration) *ProtoClient {
	return &ProtoClient{
		Conn:    natsconn,
		Timeout: timeout,
	}
}

func (c *ProtoClient) Request(subject string, args proto.Message, reply proto.Message) error {
	errWrapMessage := "failed to make a request"
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	msg, err := c.Conn.Request(subject, data, c.Timeout)
	if err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	if err = proto.Unmarshal(msg.Data, reply); err != nil {
		return errors.Wrap(err, errWrapMessage)
	}
	return nil
}
