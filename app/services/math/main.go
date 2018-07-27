package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nats-io/go-nats"

	xnats "github.com/sknv/micronats/app/lib/nats"
	xhttp "github.com/sknv/micronats/app/lib/net/http"
	"github.com/sknv/micronats/app/services/math/cfg"
	math "github.com/sknv/micronats/app/services/math/server"
)

const (
	serviceName     = "math"
	shutdownTimeout = 60 * time.Second
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

	// run the service
	natsrouter := xnats.Router{
		Conn:  natsconn,
		Queue: serviceName,
	}
	mathserver := math.NewServer(natsconn)
	mathserver.Route(&natsrouter)
	log.Printf("[INFO] %s service started", serviceName)
	defer log.Printf("[INFO] %s service stopped", serviceName)

	// run the http server
	var healthcheck healthcheck
	httprouter := http.NewServeMux()
	httprouter.HandleFunc("/healthz", healthcheck.healthz)
	xhttp.ListenAndServe(cfg.Addr, httprouter, shutdownTimeout)
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
