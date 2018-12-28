package nats

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

type ctxReq int

const (
	ctxKeyRequestID ctxReq = 0
)

func WithLimit(limit int) func(HandlerFunc) HandlerFunc {
	sem := make(chan struct{}, limit)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg *Msg) (Marshaller, error) {
			acquire()       // before request
			defer release() // after request
			return next(ctx, msg)
		}
	}
}

func WithLogger(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, msg *Msg) (out Marshaller, err error) {
		start := time.Now()

		defer func() {
			prefix := RequestIDFrom(ctx)
			statusCode := ServerHTTPStatusFromErrorCode(OK)
			if err != nil {
				status, _ := StatusFromError(err)
				statusCode = ServerHTTPStatusFromErrorCode(status.StatusCode())
				log.Printf(
					"[WARN] (%s) %s - status: %d, error: %s, in: %s",
					prefix, msg.Subject, statusCode, err, time.Since(start),
				)
				return
			}
			log.Printf(
				"[INFO] (%s) %s - status: %d, in: %s",
				prefix, msg.Subject, statusCode, time.Since(start),
			)
		}()

		out, err = next(ctx, msg)
		return
	}
}

func WithRequestID(next HandlerFunc) HandlerFunc {
	var reqid uint64
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	return func(ctx context.Context, msg *Msg) (Marshaller, error) {
		myid := atomic.AddUint64(&reqid, 1)
		requestID := fmt.Sprintf("%s-%06d", hostname, myid) // 6 digits
		ctx = context.WithValue(ctx, ctxKeyRequestID, requestID)
		return next(ctx, msg)
	}
}

func RequestIDFrom(ctx context.Context) string {
	reqID, _ := ctx.Value(ctxKeyRequestID).(string)
	return reqID
}
