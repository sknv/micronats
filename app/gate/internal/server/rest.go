package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	auth "github.com/sknv/micronats/app/auth/rpc"
	"github.com/sknv/micronats/app/gate/internal/handler"

	xhttp "github.com/sknv/micronats/app/lib/x/http"
	xnats "github.com/sknv/micronats/app/lib/x/nats"
)

func RegisterRestHandler(router chi.Router, rest *handler.Rest) {
	restServer := &restServer{
		rest:   rest,
		router: router,
	}
	restServer.route()
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type restServer struct {
	rest   *handler.Rest
	router chi.Router
}

func (s *restServer) route() {
	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/authorize", s.authorize)
		r.Post("/refresh-tokens", s.refreshTokens)
		r.Delete("/revoke-token", s.revokeToken)
	})
}

func (s *restServer) authorize(w http.ResponseWriter, r *http.Request) {
	in := new(auth.Credentials)
	s.decodeJSONInOrAbort(w, r, in)
	out, err := s.rest.Authorize(r.Context(), in)
	s.abortIfError(w, r, err)
	render.JSON(w, r, out)
}

func (s *restServer) refreshTokens(w http.ResponseWriter, r *http.Request) {
	in := new(auth.RefreshToken)
	s.decodeJSONInOrAbort(w, r, in)
	out, err := s.rest.RefreshTokens(r.Context(), in)
	s.abortIfError(w, r, err)
	render.JSON(w, r, out)
}

func (s *restServer) revokeToken(w http.ResponseWriter, r *http.Request) {
	in := new(auth.RefreshToken)
	s.decodeJSONInOrAbort(w, r, in)
	err := s.rest.RevokeToken(r.Context(), in)
	s.abortIfError(w, r, err)
	render.Status(r, http.StatusNoContent)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *restServer) decodeJSONInOrAbort(w http.ResponseWriter, r *http.Request, in interface{}) {
	if err := render.DecodeJSON(r.Body, in); err != nil {
		log.Print("[ERROR] decode json: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		xhttp.AbortHandler()
	}
}

func (s *restServer) abortIfError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	log.Print("[WARN] abort: ", err)

	// process as a xnats status
	cause := errors.Cause(err)
	status, _ := xnats.StatusFromError(cause)
	httpCode := xnats.ServerHTTPStatusFromErrorCode(status.StatusCode())
	if httpCode != http.StatusInternalServerError {
		http.Error(w, status.Message, httpCode)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	xhttp.AbortHandler()
}
