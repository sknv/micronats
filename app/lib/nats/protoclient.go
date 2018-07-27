package nats

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

const (
	errRequestMessage = "failed to make a request"
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
	data, err := proto.Marshal(args)
	if err != nil {
		return errors.Wrap(err, errRequestMessage)
	}

	msg, err := c.Conn.Request(subject, data, c.Timeout)
	if err != nil {
		return errors.Wrap(err, errRequestMessage)
	}

	if err = proto.Unmarshal(msg.Data, reply); err != nil {
		return errors.Wrap(err, errRequestMessage)
	}
	return nil
}
