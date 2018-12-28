package rpc

import (
	"context"

	xnats "github.com/sknv/micronats/app/lib/x/nats"
)

const (
	TopicAuthorize     = "/auth/authorize"
	TopicRefreshTokens = "/auth/refresh-tokens"
	TopicRevokeToken   = "/auth/revoke-token"
)

type Auth struct {
	*xnats.RemoteClient
}

func (c *Auth) Authorize(ctx context.Context, in *Credentials) (*Tokens, error) {
	out := new(Tokens)
	if err := c.Call(ctx, TopicAuthorize, in, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Auth) RefreshTokens(ctx context.Context, in *RefreshToken) (*Tokens, error) {
	out := new(Tokens)
	if err := c.Call(ctx, TopicRefreshTokens, in, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Auth) RevokeToken(ctx context.Context, in *RefreshToken) (*Empty, error) {
	out := new(Empty)
	if err := c.Call(ctx, TopicRevokeToken, in, out); err != nil {
		return nil, err
	}
	return out, nil
}
