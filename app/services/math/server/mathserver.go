package server

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/services/math/rpc"
)

const (
	mathQueue = "math"
)

func RegisterMathServer(natsServer *xnats.Server, math rpc.Math) {
	mathServer := mathServer{
		math:      math,
		publisher: xnats.NewProtoPublisher(natsServer.Conn),
	}
	mathServer.route(natsServer)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	math      rpc.Math
	publisher *xnats.ProtoPublisher
}

func (s *mathServer) route(natsServer *xnats.Server) {
	natsServer.Handle(rpc.CirclePattern, mathQueue, withLogger(s.Circle))
	natsServer.Handle(rpc.RectPattern, mathQueue, withLogger(s.Rect))
}

func (s *mathServer) Circle(ctx context.Context, message *nats.Msg) {
	args := new(rpc.CircleArgs)
	if err := proto.Unmarshal(message.Data, args); err != nil {
		panic(err) // todo: return error
	}

	reply, err := s.math.Circle(ctx, args)
	if err != nil {
		panic(err) // todo: return error
	}

	if err = s.publisher.Publish(message.Reply, reply); err != nil {
		panic(err) // todo: return error
	}

	// if err = s.TryPublish(message.Reply, reply, err); err != nil {
	// 	panic(err) // todo: return error
	// }
}

func (s *mathServer) Rect(ctx context.Context, message *nats.Msg) {
	args := new(rpc.RectArgs)
	if err := proto.Unmarshal(message.Data, args); err != nil {
		panic(err) // todo: return error
	}

	reply, err := s.math.Rect(ctx, args)
	if err != nil {
		panic(err) // todo: return error
	}

	if err = s.publisher.Publish(message.Reply, reply); err != nil {
		panic(err) // todo: return error
	}
}

// ----------------------------------------------------------------------------
// middleware example
// ----------------------------------------------------------------------------

func withLogger(next xnats.HandlerFunc) xnats.HandlerFunc {
	return func(ctx context.Context, msg *nats.Msg) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request \"%s\" processed in %s", msg.Subject, time.Since(start))
		}()
		next(ctx, msg)
	}
}
