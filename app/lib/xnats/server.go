package xnats

import (
	"context"
	"log"

	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xnats/message"
)

// HandlerFunc recievies context, subject, replyTo, message and returns an error
type HandlerFunc func(context.Context, string, string, *message.Message) error

type Server struct {
	EncConn *nats.EncodedConn
}

func NewServer(encConn *nats.EncodedConn) *Server {
	return &Server{EncConn: encConn}
}

func (s *Server) Handle(subject, queue string, handlerFn HandlerFunc) *nats.Subscription {
	sub, err := s.EncConn.QueueSubscribe(subject, queue, func(_subject, replyTo string, msg *message.Message) {
		s.handleMessageAsync(subject, replyTo, msg, handlerFn)
	})
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", "failed to set a message handler for "+subject, err)
	}
	return sub
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Server) handleMessageAsync(subject, replyTo string, message *message.Message, handlerFn HandlerFunc) {
	go func() { // process messages in a goroutine
		// recover on panic
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Printf("[PANIC] recover: %s", rvr)
			}
		}()

		// todo: fill the context with metadata
		ctx := context.Background()

		// call the handler func
		if err := handlerFn(ctx, subject, replyTo, message); err != nil {
			log.Printf("[ERROR] failed to handle the subject %s: %s", subject, err)
		}
	}()
}
