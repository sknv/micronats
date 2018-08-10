package xnats

import (
	"context"
	"log"
	"reflect"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

// HandlerFunc recieves: context, subject, replyTo, decoded object
type HandlerFunc func(context.Context, string, string, interface{})

type Server struct {
	EncConn *nats.EncodedConn
}

func NewServer(encConn *nats.EncodedConn) *Server {
	return &Server{EncConn: encConn}
}

func (s *Server) Handle(subject, queue string, handlerFn HandlerFunc) (*nats.Subscription, error) {
	sub, err := s.EncConn.QueueSubscribe(subject, queue, func(_, replyTo string, msg *Message) {
		s.handleMessageAsync(subject, replyTo, msg, handlerFn)
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to set a message handler for "+subject)
	}
	return sub, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *Server) handleMessageAsync(subject, replyTo string, message *Message, handlerFn HandlerFunc) {
	go func() { // process messages in a goroutine
		// decode the message
		decodedObj, err := s.decodeHandlerFuncObject(subject, message, handlerFn)
		if err != nil {
			log.Print("[ERROR] failed to decode the message: ", err)
			return
		}

		// todo: fill the context with metadata
		ctx := context.Background()
		handlerFn(ctx, subject, replyTo, decodedObj)
	}()
}

func (s *Server) decodeHandlerFuncObject(subject string, message *Message, handlerFn HandlerFunc) (interface{}, error) {
	objType := handlerFuncObjectType(handlerFn)

	var objPtr reflect.Value
	if objType.Kind() == reflect.Ptr {
		objPtr = reflect.New(objType.Elem())
	} else {
		objPtr = reflect.New(objType)
	}

	if err := s.EncConn.Enc.Decode(subject, message.Body, objPtr.Interface()); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal the message body")
	}
	if objType.Kind() != reflect.Ptr {
		objPtr = reflect.Indirect(objPtr)
	}
	return objPtr, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func handlerFuncObjectType(handlerFn HandlerFunc) reflect.Type {
	handlerType := reflect.TypeOf(handlerFn)
	argsNum := handlerType.NumIn()
	return handlerType.In(argsNum - 1) // object must be the last argument
}
