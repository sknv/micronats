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

func (s *Server) Handle(subject, queue string, handler HandlerFunc) (*nats.Subscription, error) {
	sub, err := s.Conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		go handler(context.Background(), msg) // process messages in a goroutine
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to set a message handler")
	}
	return sub, nil
}
