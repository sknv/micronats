package main

import (
	"log"
	"os"
	"time"

	"github.com/hashicorp/logutils"

	"github.com/sknv/micronats/app/auth/config"
	"github.com/sknv/micronats/app/auth/internal/handler"
	"github.com/sknv/micronats/app/auth/internal/server"

	xjwt "github.com/sknv/micronats/app/lib/x/jwt"
	xnats "github.com/sknv/micronats/app/lib/x/nats"
	xos "github.com/sknv/micronats/app/lib/x/os"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second
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

	// connect to nats and set draining options
	nconn, err := xnats.Connect(cfg.NatsAddr, serverShutdownTimeout)
	xos.FailOnError(err, "failed connecting to nats")
	defer nconn.Drain() // schedule draining

	// create a nats server and handle nats requests
	nmux := xnats.NewMux(
		nconn.Conn,
		&xnats.Responder{Conn: nconn.Conn},
		xnats.WithRequestID,
		xnats.WithLogger,
		xnats.WithLimit(concurrentRequestLimit),
	)
	server.RegisterAuthHandler(nmux, &handler.Auth{
		JWT: &xjwt.JWT{Secret: []byte(cfg.Secret)}, // inject a JWT manager
	})

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}
