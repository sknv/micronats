package rpc

import (
	"context"
)

type Math interface {
	Circle(context.Context, *CircleArgs) (*CircleReply, error)
	Rect(context.Context, *RectArgs) (*RectReply, error)
}
