package main

import (
	"log"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/interceptors"
	"github.com/sknv/micronats/app/lib/xos"
	"github.com/sknv/micronats/app/math/cfg"
	"github.com/sknv/micronats/app/math/server"
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
	server.RegisterMathServer(natsServer, &server.MathImpl{})

	log.Print("[INFO] math service started")
	defer log.Print("[INFO] math service stopped")

	// wait for a program exit to stop the nats server
	xos.WaitForExit()
}
