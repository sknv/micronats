package nats

import (
	"log"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
)

type DrainableConn struct {
	*nats.Conn

	wg *sync.WaitGroup
}

func Connect(url string, drainTimeout time.Duration, options ...nats.Option) (*DrainableConn, error) {
	log.Printf("[INFO] connecting to nats on %s...", url)

	// create a waitgroup required for nats shutdown
	wg := new(sync.WaitGroup)
	wg.Add(1)

	newOpts := append(
		[]nats.Option{
			nats.DrainTimeout(drainTimeout),
			nats.ClosedHandler(func(_ *nats.Conn) { wg.Done() }),
		},
		options...,
	)
	conn, err := nats.Connect(url, newOpts...)
	if err != nil {
		return nil, err
	}

	return &DrainableConn{
		Conn: conn,
		wg:   wg,
	}, nil
}

// Drain closes the nats connection.
// Timeout should be provided to "connect" function as an option: nats.DrainTimeout(shutdownTimeout)
func (c *DrainableConn) Drain() {
	if err := c.Conn.Drain(); err != nil {
		log.Print("[ERROR] failed to close nats connection: ", err)
		return
	}
	c.wg.Wait()
	log.Print("[INFO] nats connection closed")
}
