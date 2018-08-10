package server

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xnats"
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
	encConn   *nats.EncodedConn
	math      rpc.Math
	publisher *xnats.Publisher
}

func newMathServer(encConn *nats.EncodedConn, math rpc.Math) *mathServer {
	return &mathServer{
		encConn:   encConn,
		math:      math,
		publisher: xnats.NewPublisher(encConn),
	}
}

// map a request to a pattern
func (s *mathServer) route(natsServer *xnats.Server) {
	natsServer.Handle(rpc.CircleSubject, mathQueue, withLogger(s.circle))
	natsServer.Handle(rpc.RectSubject, mathQueue, withLogger(s.rect))
}

func (s *mathServer) circle(ctx context.Context, subject, replyTo string, message *xnats.Message) {
	args := new(rpc.CircleArgs)
	if err := s.encConn.Enc.Decode(subject, message.Body, args); err != nil {
		log.Print("[ERROR] failed to decode the message body: ", err)
		stat := xnats.ErrorStatus(xnats.StatusInvalidArgument, err.Error())
		if err = s.publisher.Publish(replyTo, nil, stat); err != nil {
			log.Print("[ERROR] failed to publish the error status: ", err)
		}
		return
	}

	reply, err := s.math.Circle(ctx, args)
	if err = s.publisher.Publish(replyTo, reply, err); err != nil {
		log.Print("[ERROR] failed to publish the reply: ", err)
	}
}

func (s *mathServer) rect(ctx context.Context, subject, replyTo string, message *xnats.Message) {
	args := new(rpc.RectArgs)
	if err := s.encConn.Enc.Decode(subject, message.Body, args); err != nil {
		log.Print("[ERROR] failed to decode the message body: ", err)
		stat := xnats.ErrorStatus(xnats.StatusInvalidArgument, err.Error())
		if err = s.publisher.Publish(replyTo, nil, stat); err != nil {
			log.Print("[ERROR] failed to publish the error status: ", err)
		}
		return
	}

	reply, err := s.math.Rect(ctx, args)
	if err = s.publisher.Publish(replyTo, reply, err); err != nil {
		log.Print("[ERROR] failed to publish the reply: ", err)
	}
}

// ----------------------------------------------------------------------------
// middleware example
// ----------------------------------------------------------------------------

func withLogger(next xnats.HandlerFunc) xnats.HandlerFunc {
	fn := func(ctx context.Context, subject, replyTo string, msg *xnats.Message) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request %s processed in %s", subject, time.Since(start))
		}()
		next(ctx, subject, replyTo, msg)
	}
	return fn
}
