package xnats

import (
	"context"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type HandlerFunc func(context.Context, *nats.Msg)

type Server struct {
	Conn *nats.Conn
}

func NewServer(conn *nats.Conn) *Server {
	return &Server{Conn: conn}
}

func (s *Server) Handle(subject, queue string, handler HandlerFunc) (*nats.Subscription, error) {
	sub, err := s.Conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		s.handleMessage(msg, handler)
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to set a message handler for "+subject)
	}
	return sub, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Server) handleMessage(message *nats.Msg, handler HandlerFunc) {
	ctx := context.Background()
	go handler(ctx, message) // process messages in a goroutine
}
