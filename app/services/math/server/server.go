package server

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"

	xnats "github.com/sknv/micronats/app/lib/nats"
	xmath "github.com/sknv/micronats/app/services/math/service"
)

type Server struct {
	*xnats.ProtoPublisher
}

func NewServer(natsconn *nats.Conn) *Server {
	return &Server{ProtoPublisher: xnats.NewProtoPublisher(natsconn)}
}

func (s *Server) Route(router *xnats.Router) {
	router.Handle("/math/rect", withLogger(s.HandleRect))
	router.Handle("/math/circle", withLogger(s.HandleCircle))
}

func (s *Server) HandleRect(ctx context.Context, message *nats.Msg) {
	var args xmath.RectArgs
	if err := proto.Unmarshal(message.Data, &args); err != nil {
		panic(err) // TODO: return error
	}

	reply, err := s.Rect(ctx, &args)
	if err = s.TryPublish(message.Reply, reply, err); err != nil {
		panic(err) // TODO: return error
	}
}

func (s *Server) HandleCircle(ctx context.Context, message *nats.Msg) {
	var args xmath.CircleArgs
	if err := proto.Unmarshal(message.Data, &args); err != nil {
		panic(err) // TODO: return error
	}

	reply, err := s.Circle(ctx, &args)
	if err = s.TryPublish(message.Reply, reply, err); err != nil {
		panic(err) // TODO: return error
	}
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Server) Rect(_ context.Context, args *xmath.RectArgs) (*xmath.RectReply, error) {
	perimeter := 2*args.Width + 2*args.Height
	square := args.Width * args.Height
	return &xmath.RectReply{
		Perimeter: perimeter,
		Square:    square,
	}, nil
}

func (s *Server) Circle(_ context.Context, args *xmath.CircleArgs) (*xmath.CircleReply, error) {
	len := 2 * math.Pi * args.Radius
	square := math.Pi * args.Radius * args.Radius
	return &xmath.CircleReply{
		Length: len,
		Square: square,
	}, nil
}

// ----------------------------------------------------------------------------
// middleware example
// ----------------------------------------------------------------------------

func withLogger(next xnats.HandlerFunc) xnats.HandlerFunc {
	fn := func(ctx context.Context, msg *nats.Msg) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request \"%s\" processed in %s", msg.Subject, time.Since(start))
		}()
		next(ctx, msg)
	}
	return fn
}
