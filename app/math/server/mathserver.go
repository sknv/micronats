package server

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/math/rpc"
)

const (
	mathQueue = "math"
)

func RegisterMathServer(natsServer *xnats.Server, math rpc.Math) {
	mathServer := newMathServer(natsServer.Conn, math)
	mathServer.route(natsServer)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	math      rpc.Math
	publisher *xnats.ProtoPublisher
}

func newMathServer(natsConn *nats.Conn, math rpc.Math) *mathServer {
	return &mathServer{
		math:      math,
		publisher: xnats.NewProtoPublisher(natsConn),
	}
}

// map a request to a pattern
func (s *mathServer) route(natsServer *xnats.Server) {
	natsServer.Handle(rpc.CirclePattern, mathQueue, withLogger(s.circle))
	natsServer.Handle(rpc.RectPattern, mathQueue, withLogger(s.rect))
}

func (s *mathServer) circle(ctx context.Context, message *nats.Msg) {
	args := new(rpc.CircleArgs)
	if err := proto.Unmarshal(message.Data, args); err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
	}

	reply, err := s.math.Circle(ctx, args)
	if err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
	}

	if err = s.publisher.Publish(message.Reply, reply); err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
	}
}

func (s *mathServer) rect(ctx context.Context, message *nats.Msg) {
	args := new(rpc.RectArgs)
	if err := proto.Unmarshal(message.Data, args); err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
	}

	reply, err := s.math.Rect(ctx, args)
	if err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
	}

	if err = s.publisher.Publish(message.Reply, reply); err != nil {
		log.Print("[ERROR] ", err) // todo: return error
		return
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
