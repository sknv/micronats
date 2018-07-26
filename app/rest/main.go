package main

import (
	"time"

	"github.com/go-chi/chi"

	xchi "github.com/sknv/micronats/app/lib/chi"
	xhttp "github.com/sknv/micronats/app/lib/net/http"
	"github.com/sknv/micronats/app/rest/cfg"
)

const (
	concurrentRequestLimit = 1000
	requestTimeout         = 60 * time.Second
	shutdownTimeout        = 60 * time.Second
)

func main() {
	cfg := cfg.Parse()

	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)
	xchi.UseTimeout(router, requestTimeout)

	xhttp.ListenAndServe(cfg.Addr, router, shutdownTimeout)
}
