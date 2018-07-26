package chi

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	lib "github.com/sknv/micronats/app/lib/middleware"
)

func UseDefaultMiddleware(router chi.Router) {
	router.Use(
		middleware.RealIP, middleware.Logger, middleware.Recoverer, lib.Recoverer,
	)
}

func UseThrottle(router chi.Router, concurrentRequestLimit int) {
	router.Use(middleware.Throttle(concurrentRequestLimit))
}

func UseTimeout(router chi.Router, requestTimeout time.Duration) {
	router.Use(middleware.Timeout(requestTimeout))
}

func UseLimitHandler(router chi.Router, requestLimit float64) {
	limiter := tollbooth.NewLimiter(requestLimit, nil)
	router.Use(tollbooth_chi.LimitHandler(limiter))
}

func WithLimitHandler(router chi.Router, requestLimit float64) chi.Router {
	limiter := tollbooth.NewLimiter(requestLimit, nil)
	return router.With(tollbooth_chi.LimitHandler(limiter))
}
