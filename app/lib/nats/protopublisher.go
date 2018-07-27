package nats

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

const (
	errPublishMessage = "failed to publish a reply"
)

type ProtoPublisher struct {
	Conn *nats.Conn
}

func NewProtoPublisher(natsconn *nats.Conn) *ProtoPublisher {
	return &ProtoPublisher{Conn: natsconn}
}

func (p *ProtoPublisher) Publish(subject string, reply proto.Message) error {
	data, err := proto.Marshal(reply)
	if err != nil {
		return errors.Wrap(err, errPublishMessage)
	}

	if err = p.Conn.Publish(subject, data); err != nil {
		return errors.Wrap(err, errPublishMessage)
	}
	return nil
}

func (p *ProtoPublisher) TryPublish(subject string, reply proto.Message, err error) error {
	if err != nil {
		return errors.Wrap(err, errPublishMessage)
	}
	return p.Publish(subject, reply)
}
