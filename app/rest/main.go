package main

import (
	"log"
	"time"

	"github.com/go-chi/chi"
	consul "github.com/hashicorp/consul/api"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"

	"github.com/sknv/micronats/app/lib/xchi"
	"github.com/sknv/micronats/app/lib/xconsul"
	"github.com/sknv/micronats/app/lib/xfabio"
	"github.com/sknv/micronats/app/lib/xhttp"
	"github.com/sknv/micronats/app/lib/xos"
	"github.com/sknv/micronats/app/rest/cfg"
	"github.com/sknv/micronats/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second

	serviceName         = "rest"
	healthCheckProtocol = "http://"
	healthCheckURL      = "/healthz"
	healthCheckInterval = "10s"
	healthCheckTimeout  = "1s"
)

func main() {
	cfg := cfg.Parse()

	// connect to NATS
	natsConn, _ := nats.Connect(cfg.NatsAddr)
	encConn, err := nats.NewEncodedConn(natsConn, protobuf.PROTOBUF_ENCODER)
	xos.FailOnError(err, "failed to connect to NATS")
	defer encConn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle requests
	rest := server.NewRestServer(encConn)
	rest.Route(router)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get(healthCheckURL, health.Check)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// register current service in consul and schedule a deregistration
	consulClient := registerConsulService(cfg)
	defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}

// ----------------------------------------------------------------------------
// consul section
// ----------------------------------------------------------------------------

func registerConsulService(config *cfg.Config) *xconsul.Client {
	consulClient, err := xconsul.NewClient(config.ConsulAddr)
	if err != nil {
		log.Print("[ERROR] failed to connect to consul: ", err)
		return nil
	}

	healthCheck := &consul.AgentServiceCheck{
		Name:     "rest api health check",
		HTTP:     healthCheckProtocol + config.Addr + healthCheckURL,
		Interval: healthCheckInterval,
		Timeout:  healthCheckTimeout,
	}
	if err = consulClient.RegisterCurrentService(
		config.Addr, serviceName, xfabio.Tags(serviceName), consul.AgentServiceChecks{healthCheck},
	); err != nil {
		log.Print("[ERROR] failed to register current service: ", err)
		return nil
	}
	return consulClient
}

func deregisterConsulService(consulClient *xconsul.Client) {
	if consulClient == nil {
		return
	}

	if err := consulClient.DeregisterCurrentService(); err != nil {
		log.Print("[ERROR] failed to deregister current service: ", err)
	}
}
