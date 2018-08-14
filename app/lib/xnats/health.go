package xnats

import (
	"context"
	"net/http"
	"time"

	"github.com/nats-io/go-nats"
)

const (
	healthCheckSubject = "/healthz"
)

func RegisterHealthServer(conn *nats.Conn) {
	healthServer := &healthServer{conn: conn}
	conn.QueueSubscribe(healthCheckSubject, "", healthServer.check)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type healthServer struct {
	conn *nats.Conn
}

func (s *healthServer) check(message *nats.Msg) {
	s.conn.Publish(message.Reply, []byte("ok"))
}

// ----------------------------------------------------------------------------
// http 1.1 section
// ----------------------------------------------------------------------------

type HealthServer struct {
	Conn         *nats.Conn
	CheckTimeout time.Duration
}

func NewHealthServer(conn *nats.Conn, checkTimeout time.Duration) *HealthServer {
	return &HealthServer{
		Conn:         conn,
		CheckTimeout: checkTimeout,
	}
}

func (s *HealthServer) Check(w http.ResponseWriter, _ *http.Request) {
	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.CheckTimeout)
	defer cancel()

	_, err := s.Conn.RequestWithContext(ctx, healthCheckSubject, []byte("ping"))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("service unavailable"))
		return
	}
	w.Write([]byte("ok"))
}
