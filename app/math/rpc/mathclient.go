package rpc

import (
	"context"

	"github.com/nats-io/go-nats"

	"github.com/sknv/micronats/app/lib/xnats"
)

type MathClient struct {
	*xnats.RemoteClient
}

func NewMathClient(encConn *nats.EncodedConn) Math {
	return &MathClient{RemoteClient: xnats.NewRemoteClient(encConn)}
}

func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
	reply := new(CircleReply)
	if err := c.Call(ctx, CircleSubject, args, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *MathClient) Rect(ctx context.Context, args *RectArgs) (*RectReply, error) {
	reply := new(RectReply)
	if err := c.Call(ctx, RectSubject, args, reply); err != nil {
		return nil, err
	}
	return reply, nil
}
