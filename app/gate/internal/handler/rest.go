package handler

import (
	"context"

	auth "github.com/sknv/micronats/app/auth/rpc"
	"github.com/sknv/micronats/app/core"
)

type Rest struct {
	Broker *core.Broker
}

func (r *Rest) Authorize(ctx context.Context, in *auth.Credentials) (*auth.Tokens, error) {
	newCtx, cancel := context.WithTimeout(ctx, core.RPCTimeout)
	defer cancel()

	return r.Broker.Auth.Authorize(newCtx, in)
}

func (r *Rest) RefreshTokens(ctx context.Context, in *auth.RefreshToken) (*auth.Tokens, error) {
	newCtx, cancel := context.WithTimeout(ctx, core.RPCTimeout)
	defer cancel()

	return r.Broker.Auth.RefreshTokens(newCtx, in)
}

func (r *Rest) RevokeToken(ctx context.Context, in *auth.RefreshToken) error {
	newCtx, cancel := context.WithTimeout(ctx, core.RPCTimeout)
	defer cancel()

	_, err := r.Broker.Auth.RevokeToken(newCtx, in)
	return err
}
