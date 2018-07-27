package nats

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type ProtoPublisher struct {
	Conn *nats.Conn
}

func NewProtoPublisher(natsconn *nats.Conn) *ProtoPublisher {
	return &ProtoPublisher{Conn: natsconn}
}

func (p *ProtoPublisher) Publish(subject string, reply proto.Message) error {
	errWrapMessage := "failed to publish a reply"
	data, err := proto.Marshal(reply)
	if err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	if err = p.Conn.Publish(subject, data); err != nil {
		return errors.Wrap(err, errWrapMessage)
	}
	return nil
}

func (p *ProtoPublisher) TryPublish(subject string, reply proto.Message, err error) error {
	errWrapMessage := "failed to publish a reply"
	if err != nil {
		return errors.Wrap(err, errWrapMessage)
	}
	return p.Publish(subject, reply)
}
