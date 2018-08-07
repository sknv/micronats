package xnats

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type ProtoPublisher struct {
	Conn *nats.Conn
}

func NewProtoPublisher(conn *nats.Conn) *ProtoPublisher {
	return &ProtoPublisher{Conn: conn}
}

func (p *ProtoPublisher) Publish(subject string, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "failed to marshal a message to protobuf")
	}

	if err = p.Conn.Publish(subject, data); err != nil {
		return errors.Wrap(err, "failed to publish a message")
	}
	return nil
}
