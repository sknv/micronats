package interceptors

import (
	"context"
	"log"
	"time"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/message"
)

func WithLogger(next xnats.HandlerFunc) xnats.HandlerFunc {
	return func(ctx context.Context, subject, replyTo string, msg *message.Message) error {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request %s processed in %s", subject, time.Since(start))
		}()
		return next(ctx, subject, replyTo, msg)
	}
}
