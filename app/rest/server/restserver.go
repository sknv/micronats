package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xhttp"
	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/status"
	math "github.com/sknv/micronats/app/math/rpc"
)

func RegisterRestServer(encConn *nats.EncodedConn, router chi.Router) {
	restServer := newRestServer(encConn)
	restServer.route(router)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type restServer struct {
	mathClient math.Math
}

func newRestServer(encConn *nats.EncodedConn) *restServer {
	return &restServer{mathClient: math.NewMathClient(encConn)}
}

func (s *restServer) route(router chi.Router) {
	router.Route("/math", func(r chi.Router) {
		r.Get("/circle", s.circle)
		r.Get("/rect", s.rect)
	})
}

func (s *restServer) circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	// add sample metadata
	ctx := xnats.WithMetaValue(context.Background(), "foo", "bar")

	// set the reply timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	reply, err := s.mathClient.Circle(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

func (s *restServer) rect(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	width := parseFloat(w, queryParams.Get("w"))
	height := parseFloat(w, queryParams.Get("h"))
	args := math.RectArgs{
		Width:  width,
		Height: height,
	}

	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.mathClient.Rect(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func parseFloat(w http.ResponseWriter, s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Print("[ERROR] parse float: ", err)
		http.Error(w, "argument must be a float number", http.StatusBadRequest)
		xhttp.AbortHandler()
	}
	return val
}

func abortOnError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Print("[ERROR] abort on error: ", err)

	// process as a xnats status
	cause := errors.Cause(err)
	stat, _ := status.FromError(cause)
	httpCode := status.ServerHTTPStatusFromErrorCode(stat.StatusCode())
	if httpCode != http.StatusInternalServerError {
		http.Error(w, stat.Message, httpCode)
		xhttp.AbortHandler()
	}
	xhttp.AbortHandlerWithInternalError(w)
}
