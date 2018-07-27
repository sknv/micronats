package internal

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/go-nats"

	xnats "github.com/sknv/micronats/app/lib/nats"
)

type Server struct {
	NatsConn *nats.Conn
}

func (s *Server) Route(router *xnats.Router) {
	router.Handle("/math/hello", withLogger(s.Hello))
}

func (s *Server) Hello(_ context.Context, message *nats.Msg) {
	name := string(message.Data)
	log.Print("hello, " + name)
	time.Sleep(3 * time.Second)
	s.NatsConn.Publish(message.Reply, []byte("Hello, "+name))
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
