package rpc

import (
	"context"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	"github.com/sknv/micronats/app/lib/xnats/rpc"
)

const (
	CirclePattern = "/rpc/math/circle"
	RectPattern   = "/rpc/math/rect"
)

type MathClient struct {
	*rpc.ProtoClient
}

func NewMathClient(natsConn *nats.Conn) Math {
	return &MathClient{ProtoClient: rpc.NewProtoClient(natsConn)}
}

func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
	reply := new(CircleReply)
	if err := c.Call(ctx, CirclePattern, args, reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Circle")
	}
	return reply, nil
}

func (c *MathClient) Rect(ctx context.Context, args *RectArgs) (*RectReply, error) {
	reply := new(RectReply)
	if err := c.Call(ctx, RectPattern, args, reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Rect")
	}
	return reply, nil
}
