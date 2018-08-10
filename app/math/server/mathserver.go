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
	math      rpc.Math
	responder *xnats.Responder
}

func newMathServer(encConn *nats.EncodedConn, math rpc.Math) *mathServer {
	return &mathServer{
		math:      math,
		responder: xnats.NewResponder(encConn),
	}
}

// map a request to a pattern
func (s *mathServer) route(natsServer *xnats.Server) {
	natsServer.Handle(rpc.CircleSubject, mathQueue, withLogger(s.circle))
	natsServer.Handle(rpc.RectSubject, mathQueue, withLogger(s.rect))
}

func (s *mathServer) circle(ctx context.Context, _, replyTo string, args interface{}) {
	reply, err := s.math.Circle(ctx, args.(*rpc.CircleArgs))
	err = s.responder.Response(replyTo, reply, err)
	logIfError(err)
}

func (s *mathServer) rect(ctx context.Context, _, replyTo string, args interface{}) {
	reply, err := s.math.Rect(ctx, args.(*rpc.RectArgs))
	err = s.responder.Response(replyTo, reply, err)
	logIfError(err)
}

// ----------------------------------------------------------------------------
// middleware example
// ----------------------------------------------------------------------------

func withLogger(next xnats.HandlerFunc) xnats.HandlerFunc {
	fn := func(ctx context.Context, subject, replyTo string, msg interface{}) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request \"%s\" processed in %s", subject, time.Since(start))
		}()
		next(ctx, subject, replyTo, msg)
	}
	return fn
}

func logIfError(err error) {
	if err != nil {
		log.Print("[ERROR] ", err)
	}
}
