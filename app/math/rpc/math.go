package rpc

import (
	"context"
)

const (
	CircleSubject = "/rpc/math/circle"
	RectSubject   = "/rpc/math/rect"
)

type Math interface {
	Circle(context.Context, *CircleArgs) (*CircleReply, error)
	Rect(context.Context, *RectArgs) (*RectReply, error)
}
