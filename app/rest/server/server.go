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

	xhttp "github.com/sknv/micronats/app/lib/net/http"
	math "github.com/sknv/micronats/app/services/math/service"
)

type Server struct {
	mathClient *math.Client
}

func NewServer(natsconn *nats.Conn) *Server {
	return &Server{mathClient: math.NewClient(natsconn, 10*time.Second)}
}

func (s *Server) Route(router chi.Router) {
	router.Get("/math/rect", s.Rect)
	router.Get("/math/circle", s.Circle)
}

func (s *Server) Rect(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	width := parseFloat(w, queryParams.Get("w"))
	height := parseFloat(w, queryParams.Get("h"))
	args := math.RectArgs{
		Width:  width,
		Height: height,
	}

	reply, err := s.mathClient.Rect(context.Background(), &args)
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, reply)
}

func (s *Server) Circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	reply, err := s.mathClient.Circle(context.Background(), &args)
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, reply)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func parseFloat(w http.ResponseWriter, s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Print("[ERROR] parse float: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		xhttp.AbortHandler()
	}
	return val
}
