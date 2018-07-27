package nats

import (
	"context"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type (
	HandlerFunc func(context.Context, *nats.Msg)

	Router struct {
		Conn  *nats.Conn
		Queue string
	}
)

func (s *Router) Handle(subject string, handler HandlerFunc) (*nats.Subscription, error) {
	sub, err := s.Conn.QueueSubscribe(subject, s.Queue, func(msg *nats.Msg) {
		go handler(context.Background(), msg) // process a message in a goroutine
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to set a message handler")
	}
	return sub, nil
}
