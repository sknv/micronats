package service

import (
	"context"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	xnats "github.com/sknv/micronats/app/lib/nats"
)

type Client struct {
	*xnats.ProtoClient
}

func NewClient(natsconn *nats.Conn, timeout time.Duration) *Client {
	return &Client{ProtoClient: xnats.NewProtoClient(natsconn, timeout)}
}

func (c *Client) Rect(_ context.Context, args *RectArgs) (*RectReply, error) {
	var reply RectReply
	if err := c.Request("/math/rect", args, &reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Rect")
	}
	return &reply, nil
}

func (c *Client) Circle(_ context.Context, args *CircleArgs) (*CircleReply, error) {
	var reply CircleReply
	if err := c.Request("/math/circle", args, &reply); err != nil {
		return nil, errors.Wrap(err, "failed to call Math.Circle")
	}
	return &reply, nil
}
