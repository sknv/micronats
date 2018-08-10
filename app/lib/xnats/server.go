package xnats

import (
	"context"

	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xos"
)

type HandlerFunc func(context.Context, string, string, *Message)

type Server struct {
	EncConn *nats.EncodedConn
}

func NewServer(encConn *nats.EncodedConn) *Server {
	return &Server{EncConn: encConn}
}

func (s *Server) Handle(subject, queue string, handlerFn HandlerFunc) *nats.Subscription {
	sub, err := s.EncConn.QueueSubscribe(subject, queue, func(_, replyTo string, msg *Message) {
		s.handleMessageAsync(subject, replyTo, msg, handlerFn)
	})
	xos.FailOnError(err, "failed to set a message handler for "+subject)
	return sub
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Server) handleMessageAsync(subject, replyTo string, message *Message, handlerFn HandlerFunc) {
	go func() { // process messages in a goroutine
		// todo: recover on panic

		// todo: fill the context with metadata
		ctx := context.Background()

		// call the handler func
		handlerFn(ctx, subject, replyTo, message)
	}()
}
