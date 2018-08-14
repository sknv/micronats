package server

import (
	"context"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/message"
	"github.com/sknv/micronats/app/lib/xnats/status"
	"github.com/sknv/micronats/app/math/rpc"
)

const (
	mathQueue = "math"
)

func RegisterMathServer(natsServer *xnats.Server, math rpc.Math) {
	mathServer := newMathServer(natsServer.EncConn, math)
	mathServer.route(natsServer)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	encoder   nats.Encoder
	math      rpc.Math
	publisher *xnats.Publisher
}

func newMathServer(encConn *nats.EncodedConn, math rpc.Math) *mathServer {
	return &mathServer{
		encoder:   encConn.Enc,
		math:      math,
		publisher: xnats.NewPublisher(encConn),
	}
}

// map a request to a subject
func (s *mathServer) route(natsServer *xnats.Server) {
	natsServer.Handle(rpc.CircleSubject, mathQueue, s.circle)
	natsServer.Handle(rpc.RectSubject, mathQueue, s.rect)
}

func (s *mathServer) circle(ctx context.Context, subject, replyTo string, message *message.Message) error {
	args := new(rpc.CircleArgs)
	if err := s.decodeArgs(subject, replyTo, message, args); err != nil {
		return err
	}

	reply, err := s.math.Circle(ctx, args)
	if err = s.publisher.Publish(replyTo, reply, err); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}

func (s *mathServer) rect(ctx context.Context, subject, replyTo string, message *message.Message) error {
	args := new(rpc.RectArgs)
	if err := s.decodeArgs(subject, replyTo, message, args); err != nil {
		return err
	}

	reply, err := s.math.Rect(ctx, args)
	if err = s.publisher.Publish(replyTo, reply, err); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *mathServer) decodeArgs(subject, replyTo string, message *message.Message, args interface{}) error {
	err := s.encoder.Decode(subject, message.Body, args)
	if err == nil {
		return nil
	}

	err = errors.WithMessage(err, "failed to decode the message body")
	status := status.Error(status.InvalidArgument, err.Error())
	if puberr := s.publisher.Publish(replyTo, nil, status); puberr != nil {
		puberr = errors.WithMessage(puberr, "failed to publish the error status")
		err = errors.WithMessage(err, puberr.Error())
	}
	return err
}
