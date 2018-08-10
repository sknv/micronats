package main

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xchi"
	"github.com/sknv/micronats/app/lib/xhttp"
	"github.com/sknv/micronats/app/lib/xos"
	"github.com/sknv/micronats/app/rest/cfg"
	"github.com/sknv/micronats/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second
)

func main() {
	cfg := cfg.Parse()

	// connect to NATS
	natsConn, _ := nats.Connect(cfg.NatsAddr)
	encConn, err := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	xos.FailOnError(err, "failed to connect to NATS")

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle requests
	rest := server.NewRestServer(encConn)
	rest.Route(router)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get("/healthz", health.Check)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}
