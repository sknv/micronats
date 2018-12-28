package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hashicorp/logutils"
	"github.com/nats-io/go-nats"

	auth "github.com/sknv/micronats/app/auth/rpc"
	"github.com/sknv/micronats/app/core"
	"github.com/sknv/micronats/app/gate/config"
	"github.com/sknv/micronats/app/gate/internal/handler"
	"github.com/sknv/micronats/app/gate/internal/server"

	xhttp "github.com/sknv/micronats/app/lib/x/http"
	xnats "github.com/sknv/micronats/app/lib/x/nats"
	xos "github.com/sknv/micronats/app/lib/x/os"
)

const (
	concurrentRequestLimit = 1000
	serverReadTimeout      = 30 * time.Second
	serverWriteTimeout     = 30 * time.Second
	serverShutdownTimeout  = 60 * time.Second // gte than write timeout to flush all the responses
)

func main() {
	cfg := config.Parse()

	// customize the logger
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	// connect to nats
	log.Printf("[INFO] connecting to nats on %s...", cfg.NatsAddr)
	nconn, err := nats.Connect(cfg.NatsAddr)
	xos.FailOnError(err, "failed connecting to nats")
	defer func() { // schedule closing
		nconn.Close()
		log.Print("[INFO] nats connection closed")
	}()

	// create rpc clients and combine them into a broker
	nclient := &xnats.RemoteClient{Conn: nconn}
	broker := &core.Broker{
		Auth: auth.Auth{RemoteClient: nclient},
	}

	// create and config an http router
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		xhttp.Recoverer,
		middleware.Throttle(concurrentRequestLimit),
	)

	// handle rest requests
	server.RegisterRestHandler(router, &handler.Rest{
		Broker: broker,
	})

	// start the http server
	srv := &xhttp.Server{
		Server: &http.Server{
			Addr:         cfg.GateAddr,
			Handler:      router,
			ReadTimeout:  serverReadTimeout,
			WriteTimeout: serverWriteTimeout,
		},
	}
	srv.Start()
	defer srv.Shutdown(serverShutdownTimeout) // schedule a graceful shutdown

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}
