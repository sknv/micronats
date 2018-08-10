package xnats

// import (
// 	"context"
// 	"log"
// 	"reflect"

// 	"github.com/nats-io/go-nats"
// 	"github.com/pkg/errors"

// 	"github.com/sknv/micronats/app/lib/xos"
// )

// // HandlerFunc recieves: context, subject, replyTo, args object
// type HandlerFunc interface{}

// type Server struct {
// 	EncConn *nats.EncodedConn
// }

// func NewServer(encConn *nats.EncodedConn) *Server {
// 	return &Server{EncConn: encConn}
// }

// func (s *Server) Handle(subject, queue string, handlerFn HandlerFunc) *nats.Subscription {
// 	objType := handlerFuncObjectType(handlerFn) // validate the handler func and get the args type
// 	handlerFnVal := reflect.ValueOf(handlerFn)  // get the handler func itself

// 	sub, err := s.EncConn.QueueSubscribe(subject, queue, func(_, replyTo string, msg *Message) {
// 		s.handleMessageAsync(subject, replyTo, msg, objType, handlerFnVal)
// 	})
// 	xos.FailOnError(err, "failed to set a message handler for "+subject)
// 	return sub
// }

// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------

// func (s *Server) handleMessageAsync(subject, replyTo string, message *Message, objType reflect.Type, handlerFnVal reflect.Value) {
// 	go func() { // process messages in a goroutine
// 		// todo: recover on panic

// 		// decode the message and fill the args object
// 		objVal, err := s.decodeHandlerFuncObject(subject, message, objType)
// 		if err != nil {
// 			log.Print("[ERROR] failed to decode the message: ", err)
// 			return
// 		}

// 		// todo: fill the context with metadata
// 		ctx := context.Background()

// 		// call the handler func
// 		ctxVal := reflect.ValueOf(ctx)
// 		subjectVal := reflect.ValueOf(subject)
// 		replyToVal := reflect.ValueOf(replyTo)
// 		handlerArgs := []reflect.Value{ctxVal, subjectVal, replyToVal, objVal}
// 		handlerFnVal.Call(handlerArgs)
// 	}()
// }

// func (s *Server) decodeHandlerFuncObject(subject string, message *Message, objType reflect.Type) (reflect.Value, error) {
// 	objVal, err := s.decodeObject(subject, message.Body, objType)
// 	if err != nil {
// 		return reflect.Value{}, errors.WithMessage(err, "failed to decode the message body")
// 	}
// 	return objVal, nil
// }

// func (s *Server) decodeObject(subject string, data []byte, objType reflect.Type) (reflect.Value, error) {
// 	var objVal reflect.Value
// 	if objType.Kind() == reflect.Ptr {
// 		objVal = reflect.New(objType.Elem())
// 	} else {
// 		objVal = reflect.New(objType)
// 	}

// 	if err := s.EncConn.Enc.Decode(subject, data, objVal.Interface()); err != nil {
// 		return reflect.Value{}, errors.WithMessage(err, "failed to decode the object")
// 	}
// 	return reflect.Indirect(objVal), nil
// }

// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------

// func handlerFuncObjectType(handlerFn HandlerFunc) reflect.Type {
// 	handlerType := reflect.TypeOf(handlerFn)
// 	if handlerType.Kind() != reflect.Func {
// 		panic("xnats: handler needs to be a func")
// 	}
// 	argsNum := handlerType.NumIn()
// 	if argsNum != 4 {
// 		panic("xnats: handler func must recieve 4 args: context, string, string, object")
// 	}
// 	return handlerType.In(argsNum - 1) // args object must be the last argument
// }
