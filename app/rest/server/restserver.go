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
	"github.com/sknv/micronats/app/lib/xnats/status"
	math "github.com/sknv/micronats/app/math/rpc"
)

type RestServer struct {
	mathClient math.Math
}

func NewRestServer(natsConn *nats.Conn) *RestServer {
	return &RestServer{mathClient: math.NewMathClient(natsConn)}
}

func (s *RestServer) Route(router chi.Router) {
	router.Route("/math", func(r chi.Router) {
		r.Get("/circle", s.circle)
		r.Get("/rect", s.rect)
	})
}

func (s *RestServer) circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.mathClient.Circle(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

func (s *RestServer) rect(w http.ResponseWriter, r *http.Request) {
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

	// process as an xnats.Status
	cause := errors.Cause(err)
	stat, _ := status.FromError(cause)
	httpCode := status.ServerHTTPStatusFromErrorCode(stat.StatusCode())
	if httpCode != http.StatusInternalServerError {
		http.Error(w, stat.GetMessage(), httpCode)
		xhttp.AbortHandler()
	}
	xhttp.AbortHandlerWithInternalError(w)
}
