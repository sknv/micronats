package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/nats-io/go-nats"

	xchi "github.com/sknv/micronats/app/lib/chi"
	xhttp "github.com/sknv/micronats/app/lib/net/http"
	"github.com/sknv/micronats/app/rest/cfg"
	"github.com/sknv/micronats/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	shutdownTimeout        = 60 * time.Second
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type healthcheck struct{}

func (*healthcheck) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func main() {
	cfg := cfg.Parse()

	// connect to the NATS server
	natsconn, err := nats.Connect(cfg.NatsURL)
	failOnError(err, "failed to connect to NATS server")
	defer natsconn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// route the server
	srv := server.NewServer(natsconn)
	srv.Route(router)

	// run the http server
	var healthcheck healthcheck
	router.Get("/healthz", healthcheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, router, shutdownTimeout)
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
