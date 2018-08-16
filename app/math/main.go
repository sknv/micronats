package main

import (
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	consul "github.com/hashicorp/consul/api"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"

	"github.com/sknv/micronats/app/lib/xconsul"
	"github.com/sknv/micronats/app/lib/xhttp"
	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/interceptors"
	"github.com/sknv/micronats/app/lib/xos"
	"github.com/sknv/micronats/app/math/cfg"
	math "github.com/sknv/micronats/app/math/server"
)

const (
	healthServerShutdownTimeout = 60 * time.Second

	serviceName         = "math"
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

	// handle nats requests
	natsServer := xnats.NewServer(encConn, interceptors.WithLogger)
	math.RegisterMathServer(natsServer, &math.MathImpl{})
	xnats.RegisterHealthServer(natsConn) // handle nats health check requests

	// config the http router for health checks
	router := chi.NewRouter()
	router.Use(middleware.RealIP, middleware.Logger)

	// handle health check requests
	healthCheckDuration, err := time.ParseDuration(healthCheckTimeout)
	xos.FailOnError(err, "failed to parse health check timeout")
	health := xnats.NewHealthServer(natsConn, healthCheckDuration)
	router.Get(healthCheckURL, health.Check)

	// start the http health check server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(healthServerShutdownTimeout)

	// register current service in consul and schedule a deregistration
	consulClient := registerConsulService(cfg)
	defer deregisterConsulService(consulClient)

	// wait for a program exit to stop the nats server
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
		Name:     "math service health check",
		HTTP:     healthCheckProtocol + config.Addr + healthCheckURL,
		Interval: healthCheckInterval,
		Timeout:  healthCheckTimeout,
	}
	if err = consulClient.RegisterCurrentService(
		config.Addr, serviceName, nil, consul.AgentServiceChecks{healthCheck},
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
