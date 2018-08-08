package xnats

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats/message"
)

type ProtoPublisher struct {
	Conn *nats.Conn
}

func NewProtoPublisher(conn *nats.Conn) *ProtoPublisher {
	return &ProtoPublisher{Conn: conn}
}

func (p *ProtoPublisher) Publish(subject string, msg proto.Message) error {
	body, err := ptypes.MarshalAny(msg)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal the message to protobuf any")
	}

	protoMsg := &message.Message{Body: body}
	data, err := proto.Marshal(protoMsg)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal the message to protobuf")
	}

	if err = p.Conn.Publish(subject, data); err != nil {
		return errors.WithMessage(err, "failed to publish the message")
	}
	return nil
}
