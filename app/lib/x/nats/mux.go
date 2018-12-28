package nats

import (
	"context"
	"errors"
	"log"

	"github.com/nats-io/go-nats"
)

var (
	errSmoke = errors.New("something went wrong")
)

type (
	HandlerFunc func(context.Context, *Msg) (Marshaller, error)

	Interceptor func(HandlerFunc) HandlerFunc
)

type Mux struct {
	Conn         *nats.Conn
	Responder    *Responder
	Interceptors []Interceptor
}

func NewMux(conn *nats.Conn, responder *Responder, interceptors ...Interceptor) *Mux {
	return &Mux{
		Conn:         conn,
		Responder:    responder,
		Interceptors: interceptors,
	}
}

func (s *Mux) Handle(subject, queue string, handle HandlerFunc) (*nats.Subscription, error) {
	return s.Conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		go s.handleMessage(&Msg{Msg: msg}, handle) // process messages in a goroutine
	})
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Mux) handleMessage(msg *Msg, handle HandlerFunc) {
	// recover on panic
	defer func() {
		if rvr := recover(); rvr != nil {
			log.Print("[ERROR] panic: ", rvr)
			if err := s.Responder.Respond(msg.Reply, nil, errSmoke); err != nil {
				log.Print("[ERROR] publish smoke error: ", err)
			}
		}
	}()

	// unmarshal message payload
	payload, err := msg.Payload()
	if err != nil {
		log.Print("[ERROR] unmarshal payload: ", err)
		if err = s.Responder.Respond(msg.Reply, nil, errSmoke); err != nil {
			log.Print("[ERROR] publish smoke error: ", err)
		}
		return
	}

	// fill the context with provided metadata
	ctx := context.Background()
	if payload.Meta != nil {
		for key, val := range payload.Meta {
			ctx = ContextWithMetaValue(ctx, key, val)
		}
	}

	// execute the interceptors
	handle = s.chainInterceptors(handle)

	// call the handler func
	out, err := handle(ctx, msg)
	if err != nil {
		if err = s.Responder.Respond(msg.Reply, nil, err); err != nil {
			log.Print("[ERROR] publish error: ", err)
		}
		return
	}

	// no need to publish anything
	if out == nil {
		return
	}

	if err = s.Responder.Respond(msg.Reply, out, nil); err != nil {
		log.Print("[ERROR] publish response: ", err)
		return
	}
}

func (s *Mux) chainInterceptors(endpoint HandlerFunc) HandlerFunc {
	if len(s.Interceptors) == 0 {
		return endpoint
	}

	handler := s.Interceptors[len(s.Interceptors)-1](endpoint)
	for i := len(s.Interceptors) - 2; i >= 0; i-- {
		handler = s.Interceptors[i](handler)
	}
	return handler
}
