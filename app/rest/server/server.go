package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/go-nats"
)

type Server struct {
	NatsConn *nats.Conn
}

func (s *Server) Route(router chi.Router) {
	router.Get("/hello/{name}", s.Hello)
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	reply, err := s.NatsConn.Request("/math/hello", []byte(name), 5*time.Second)
	if err != nil {
		panic(err)
	}
	render.PlainText(w, r, string(reply.Data))
}
