package server

import (
	"context"

	"github.com/sknv/micronats/app/auth/internal/handler"
	"github.com/sknv/micronats/app/auth/rpc"

	xnats "github.com/sknv/micronats/app/lib/x/nats"
)

const (
	authQueue = "auth"
)

func RegisterAuthHandler(natsMux *xnats.Mux, auth *handler.Auth) {
	authServer := &authServer{
		auth: auth,
		mux:  natsMux,
	}
	authServer.route()
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type authServer struct {
	auth *handler.Auth
	mux  *xnats.Mux
}

func (s *authServer) route() {
	s.mux.Handle(rpc.TopicAuthorize, authQueue, s.authorize)
	s.mux.Handle(rpc.TopicRefreshTokens, authQueue, s.refreshTokens)
	s.mux.Handle(rpc.TopicRevokeToken, authQueue, s.revokeToken)
}

func (s *authServer) authorize(ctx context.Context, msg *xnats.Msg) (xnats.Marshaller, error) {
	in := new(rpc.Credentials)
	payload, _ := msg.Payload()
	if err := in.Unmarshal(payload.GetBody()); err != nil {
		return nil, err
	}
	return s.auth.Authorize(ctx, in)
}

func (s *authServer) refreshTokens(ctx context.Context, msg *xnats.Msg) (xnats.Marshaller, error) {
	in := new(rpc.RefreshToken)
	payload, _ := msg.Payload()
	if err := in.Unmarshal(payload.GetBody()); err != nil {
		return nil, err
	}
	return s.auth.RefreshTokens(ctx, in)
}

func (s *authServer) revokeToken(ctx context.Context, msg *xnats.Msg) (xnats.Marshaller, error) {
	in := new(rpc.RefreshToken)
	payload, _ := msg.Payload()
	if err := in.Unmarshal(payload.GetBody()); err != nil {
		return nil, err
	}
	return s.auth.RevokeToken(ctx, in)
}
