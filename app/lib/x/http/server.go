package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func (s *Server) Start() {
	log.Printf("[INFO] starting an http server on %s...", s.Server.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[WARN] shutting down the http server: ", err)
		}
	}()
}

func (s *Server) Shutdown(shutdownTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		log.Print("[ERROR] failed to shutdown the http server: ", err)
	}
	log.Print("[INFO] http server stopped")
}
